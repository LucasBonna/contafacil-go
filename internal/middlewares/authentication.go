package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

func Authenticate(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") || c.Request.URL.Path == "/docs" {
			c.Next()
			return
		}

		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(bearer, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		token := strings.TrimPrefix(bearer, prefix)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		queryResp, err := deps.Core.DB.GetUserAndClientByApiKey(context.Background(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		clientDetails := schemas.ClientDetails{
			User: schemas.User{
				ID:        queryResp.UserID,
				Username:  queryResp.Username,
				ApiKey:    queryResp.ApiKey,
				Role:      queryResp.UserRole,
				ClientID:  queryResp.ClientID,
				CreatedAt: queryResp.UserCreatedAt,
				UpdatedAt: queryResp.UserUpdatedAt,
				DeletedAt: queryResp.UserDeletedAt,
			},
			Client: schemas.Client{
				ID:        queryResp.ClientID,
				Name:      queryResp.ClientName,
				Cnpj:      queryResp.ClientCnpj,
				Role:      queryResp.ClientRole,
				CreatedAt: queryResp.ClientCreatedAt,
				UpdatedAt: queryResp.ClientUpdatedAt,
				DeletedAt: queryResp.ClientDeletedAt,
			},
		}

		c.Set("clientDetails", &clientDetails)

		c.Next()
	}
}

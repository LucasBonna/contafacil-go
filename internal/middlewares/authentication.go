package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/ent/user"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

func Authenticate(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/monitoring") ||
			strings.HasPrefix(path, "/swagger") ||
			strings.HasPrefix(path, "/docs") ||
			strings.HasPrefix(path, "/auth/login") {
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

		user, err := deps.Core.DB.User.
			Query().
			Where(user.APIKey(token)).
			WithClients().
			Only(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		client := user.Edges.Clients

		clientDetails := schemas.ClientDetails{
			User: schemas.User{
				ID:        user.ID,
				Username:  user.Username,
				ApiKey:    user.APIKey,
				Role:      string(user.Role),
				ClientID:  user.ClientID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			Client: schemas.Client{
				ID:        client.ID,
				Name:      client.Name,
				Cnpj:      client.Cnpj,
				Role:      string(client.Role),
				CreatedAt: client.CreatedAt,
				UpdatedAt: client.UpdatedAt,
				DeletedAt: client.DeletedAt,
			},
		}

		c.Set("clientDetails", &clientDetails)

		c.Next()
	}
}

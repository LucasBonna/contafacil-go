package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/lucasbonna/contafacil_api/ent/user"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type jwtClaims struct {
	User                 schemas.JWTToken `json:"user"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthRouter(r *gin.Engine, deps *app.Dependencies) {
	r.POST("teste", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginHandler(deps))
	}
}

func loginHandler(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input loginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		// Find user by username
		u, err := deps.Core.DB.User.
			Query().
			Where(user.Username(input.Username)).
			WithClients().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		client := schemas.Client{
			ID:        u.ClientID,
			Name:      u.Edges.Clients.Name,
			Cnpj:      u.Edges.Clients.Cnpj,
			Role:      string(u.Edges.Clients.Role),
			CreatedAt: u.Edges.Clients.CreatedAt,
			UpdatedAt: u.Edges.Clients.UpdatedAt,
			DeletedAt: u.Edges.Clients.DeletedAt,
		}

		claims := jwtClaims{
			User: schemas.JWTToken{
				ID:       u.ID,
				Username: u.Username,
				ApiKey:   u.APIKey,
				Role:     string(u.Role),
				Client:   client,
			},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.New().String(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(config.Env.JWTSecret))
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(200, gin.H{
			"token": signedToken,
		})
	}
}

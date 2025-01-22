package server

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/middlewares"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
)

type Server struct {
	deps *app.Dependencies
}

func NewServer(dbConnStr string, rabbit *rabbitmq.RabbitMQ, deps *app.Dependencies) *Server {
	return &Server{
		deps: deps,
	}
}

func (s *Server) StartServer() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "POST", "GET", "HEAD"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
	}))

	r.Use(gin.Recovery())

	// expectedHost := config.FrontEndUrl
	expectedHost := "localhost:8000"

	r.Use(func(c *gin.Context) {
		if c.Request.Host != expectedHost {
			log.Println("request host", c.Request.Host)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	r.Use(middlewares.Authenticate(s.deps))

	r.Use(middlewares.Logger(s.deps))

	InitRouters(r, s.deps)

	r.Run(":8000")
}

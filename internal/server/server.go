package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/middlewares"
)

type Server struct {
	deps *app.Dependencies
}

func NewServer(deps *app.Dependencies) *Server {
	return &Server{
		deps: deps,
	}
}

func (s *Server) StartServer() {
	r := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Env.FrontEndUrl}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(middlewares.Authenticate(s.deps))

	r.Use(middlewares.Logger(s.deps))

	InitRouters(r, s.deps)

	r.Run(":8000")
}

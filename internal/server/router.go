package server

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/server/routes"
)

func InitRouters(r *gin.Engine, deps *app.Dependencies) {
	// Aqui vem os routers
	routes.DefaultRouter(r)

	routes.FileRouter(r, deps)

	routes.EmissionRouter(r, deps)

	routes.TestRoutes(r, deps)
}

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/server/routes"
)

func InitRouters(r *gin.Engine, deps *app.Dependencies) {
	routes.DefaultRouter(r)

	routes.AuthRouter(r, deps)

	routes.ClientRoutes(r, &deps.Core)

	routes.UserRouters(r, &deps.Core)

	routes.FileRouter(r, deps)

	routes.EmissionRouter(r, deps)

	routes.SSERouter(r, deps)

	routes.TestRoutes(r, deps)
}

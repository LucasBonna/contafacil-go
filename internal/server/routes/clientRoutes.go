package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func ClientRoutes(r *gin.Engine, core *app.CoreDependencies) {
	clientsHandler := handlers.NewClientHandlers(core)

	clients := r.Group("/clients")
	{
		clients.GET("/", clientsHandler.ListAllClients())
		clients.GET("/:id", clientsHandler.GetClient())
		clients.POST("/", clientsHandler.CreateClient())
		clients.PATCH("/:id", clientsHandler.UpdateClient())
		clients.DELETE("/:id", clientsHandler.DeleteClient())
	}
}

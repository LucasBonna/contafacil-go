package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func SSERouter(r *gin.Engine, deps *app.Dependencies) {
	sseHandler := handlers.NewSSEHandler(deps)
	r.GET("/sse/:userId", sseHandler.HandleSSE)
}

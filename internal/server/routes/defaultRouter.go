package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func DefaultRouter(r *gin.Engine) {
	r.GET("/health", handlers.HealthHandler)
}

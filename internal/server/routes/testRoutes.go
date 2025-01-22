package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func TestRoutes(r *gin.Engine, deps *app.Dependencies) {

  r.GET("/test/queue", handlers.HandlerTestQueue(deps))

  r.GET("/test/requeue", handlers.HandlerRequeueExceptions(deps))

}


package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
)

func HandlerTestQueue(deps *app.Dependencies) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("entrou handler")

		ctx.JSON(http.StatusOK, gin.H{"message": "Task enqueued successfully"})
	}
}

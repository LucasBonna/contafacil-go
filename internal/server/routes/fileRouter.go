package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/handlers"
)

func FileRouter(r *gin.Engine, deps *app.Dependencies) {
	file := r.Group("/file")
	{
		file.POST("/upload", handlers.HandlerUploadFile(deps))
		file.POST("/download/batch", handlers.HandlerDownloadBatch(deps))
		file.POST("/download/:fileId", handlers.HandlerDownloadFile(deps))
	}
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

func HandlerTestQueue(deps *app.Dependencies) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("entrou handler")
		payload := map[string]interface{}{
			"filename": "test.txt",
			"size":     1024,
		}

		err := deps.Core.QH.EnqueueTask(utils.TaskUploadFile, payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Task enqueued successfully"})
	}
}

func HandlerRequeueExceptions(deps *app.Dependencies) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deliveries, err := deps.Core.Rabbit.Consume("ExceptionQueue", "")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var reprocessedCount int
		for msg := range deliveries {
			var task utils.TaskWithMetadata
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Printf("invalid message format in exception queue: %v", err)
				msg.Reject(false)
				continue
			}

			err = deps.Core.Rabbit.Publish(task.Metadata.OriginalQueue, msg.Body)
			if err != nil {
				log.Printf("failed to requeue message to %s: %v", task.Metadata.OriginalQueue, err)
				msg.Reject(false)
				continue
			}

			msg.Ack(false)
			reprocessedCount++
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":             "messages reprocessed successfully",
			"reprocessedCount":    reprocessedCount,
			"exceptionQueueCount": len(deliveries),
		})
	}
}

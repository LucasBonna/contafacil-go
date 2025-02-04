package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

func SSERouter(r *gin.Engine, deps *app.Dependencies) {
	r.GET("/sse/:userId", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		rawUserId := c.Param("userId")
		if rawUserId == "" {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		messageChan := make(chan schemas.SSEMessage, 10)

		userId, err := uuid.Parse(rawUserId)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		deps.Core.SSEManager.AddClient(userId, messageChan)
		defer deps.Core.SSEManager.RemoveClient(userId)

		done := c.Request.Context().Done()

		c.SSEvent("connected", time.Now().Format(time.RFC3339))
		c.Writer.Flush()

		for {
			select {
			case msg := <-messageChan:
				c.SSEvent(msg.Event, msg.Data)
				c.Writer.Flush()
			case <-done:
				return
			case <-time.After(30 * time.Second):
				c.SSEvent("heartbeat", time.Now().Format(time.RFC3339))
				c.Writer.Flush()
			}
		}
	})
}

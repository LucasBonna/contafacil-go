package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/internal/app"
)

type SSEHandler struct {
	deps *app.Dependencies
}

func NewSSEHandler(deps *app.Dependencies) *SSEHandler {
	return &SSEHandler{deps: deps}
}

func (sh *SSEHandler) HandleSSE(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Configurar headers SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Registrar conexão
	msgChan := sh.deps.Core.SSEMgr.Register(userID)
	defer sh.deps.Core.SSEMgr.Unregister(userID)

	// Heartbeat
	ticker := time.NewTicker(15 * time.Second)
	cleanupTicker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	defer cleanupTicker.Stop()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	for {
		select {
		case msg := <-msgChan:
			data, _ := json.Marshal(msg)
			c.SSEvent(msg.Event, string(data))
			c.Writer.Flush()

		case <-ticker.C:
			c.SSEvent("heartbeat", nil)
			c.Writer.Flush()

		case <-cleanupTicker.C:
			sh.deps.Core.SSEMgr.ListClients()

		case <-ctx.Done():
			return
		}
	}
}

package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/internal/app"
)

type SSEHandler struct {
	deps *app.Dependencies
}

func NewSSEHandler(deps *app.Dependencies) *SSEHandler {
	return &SSEHandler{deps: deps}
}

func (h *SSEHandler) ProcessSSEUpdate(ctx context.Context, t *asynq.Task) error {
	var payload SSEUpdatePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	// Verificar presença ativa
	if !h.deps.Core.SSEMgr.IsConnected(payload.UserID) {
		return nil // Descarta se não conectado
	}

	// Enviar mensagem diretamente
	success := h.deps.Core.SSEMgr.Send(payload.UserID, payload.Message)
	if !success {
		return fmt.Errorf("failed to send SSE message")
	}

	return nil
}

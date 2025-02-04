package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type SSEHandler struct {
	Deps *app.Dependencies
}

func NewSSEHandler(deps *app.Dependencies) *SSEHandler {
	return &SSEHandler{
		Deps: deps,
	}
}

type SSEUpdatePayload struct {
	EmissionID uuid.UUID
	Status     emission.Status
	Message    string
	ClientID   uuid.UUID
	UserID     uuid.UUID
}

func (sh *SSEHandler) ProcessSSEUpdate(ctx context.Context, t *asynq.Task) error {
	var p SSEUpdatePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	message := schemas.SSEMessage{
		Event: "emission_update",
		Data: schemas.EmissionUpdate{
			EmissionID: p.EmissionID,
			Status:     p.Status,
			Message:    p.Message,
		},
	}

	retries := 3
	for i := 0; i < retries; i++ {
		ok := sh.Deps.Core.SSEManager.SendToClient(p.UserID, message)
		if ok {
			return nil
		}
		log.Printf("Retrying send to user %s (%d/%d)", p.UserID, i+1, retries)
		time.Sleep(1 * time.Second)
	}

	log.Printf("Failed to send message to user %s after %d retries", p.UserID, retries)
	return nil
}

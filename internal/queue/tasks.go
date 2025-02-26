package queue

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/internal/sse"
)

type TaskPayload interface{}

func NewTask[T TaskPayload](taskType string, payload T) (*asynq.Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload")
	}

	return asynq.NewTask(taskType, payloadBytes, asynq.MaxRetry(3)), nil
}

type SSEUpdatePayload struct {
	UserID  uuid.UUID   `json:"user_id"`
	Message sse.Message `json:"message"`
}

const (
	TypeSSEEmissionUpdate = "sse:emission_update"
)

package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type TaskPayload interface{}

func NewTask[T TaskPayload](taskType string, payload T) (*asynq.Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload")
	}

	return asynq.NewTask(taskType, payloadBytes, asynq.MaxRetry(3)), nil
}

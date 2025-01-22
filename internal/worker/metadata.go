package worker

import (
	"github.com/google/uuid"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

type MessageMetadata struct {
  RetryCount int `json:"retryCount"`
  RetryAt int64 `json:"retryAt"`
  OriginalQueue string `json:"originalQueue"`
}

type TaskWithMetadata struct {
  Id uuid.UUID
  Type utils.TaskType `json:"type"`
  Payload interface{} `json:"payload"`
  Metadata *MessageMetadata `json:"messageMetadata,omitempty"`
}

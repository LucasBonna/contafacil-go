package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type IssueGNRETaskPayload struct {
	XmlContent    string
	ClientDetails *schemas.ClientDetails
}

func NewIssueGNRETask(xmlContent string, clientDetails *schemas.ClientDetails) (*asynq.Task, error) {
	payload, err := json.Marshal(IssueGNRETaskPayload{XmlContent: xmlContent, ClientDetails: clientDetails})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeIssueGNRE, payload, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute)), nil
}

func HandleIssueGNRETask(ctx context.Context, t *asynq.Task) error {
	var p IssueGNRETaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Processing GNRE: XmlContent: %v, ClientDetails: %v", p.XmlContent, p.ClientDetails)
	return nil
}

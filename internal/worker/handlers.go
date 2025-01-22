package worker

import (
	"context"
	"log"
)

// TaskHandler é a interface que todos os handlers devem implementar.
type TaskHandler interface {
	Handle(ctx context.Context, payload interface{}) error
}

// IssueGNREHandler processa a geração de GNREs.
type IssueGNREHandler struct{}

func (h *IssueGNREHandler) Handle(ctx context.Context, payload interface{}) error {
	log.Printf("Processing GNRE issue: %v", payload)
	// Lógica de geração de GNRE
	return nil
}

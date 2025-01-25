package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type GNREHandler struct {
	Deps *app.Dependencies
}

func NewGNREHandler(deps *app.Dependencies) *GNREHandler {
	return &GNREHandler{
		Deps: deps,
	}
}

type IssueGNRETaskPayload struct {
	EmissionId    uuid.UUID
	XmlContent    string
	ClientDetails *schemas.ClientDetails
}

func (gh *GNREHandler) ProcessIssueGNRE(ctx context.Context, t *asynq.Task) error {
	var p IssueGNRETaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing GNRE: Client: %v", p.ClientDetails.Client.Name)

	log.Println("emissionId", p.EmissionId)

	tecnospeedResponse, err := gh.Deps.External.TecnospeedService.IssueGNRE(p.XmlContent, "ContaFacil", p.ClientDetails.Client.Cnpj)
	if err != nil {
		_, err := gh.Deps.Core.DB.UpdateEmissionAndGNRE(ctx, database.UpdateEmissionAndGNREParams{
			ID:      pgtype.UUID{Bytes: p.EmissionId},
			Status:  pgtype.Text{String: "FAILED", Valid: true},
			Message: pgtype.Text{String: tecnospeedResponse.Failure.Message, Valid: true},
		})
		return fmt.Errorf("issueGNRE failed: %v", err, asynq.SkipRetry)
	}

	updatedGNRE, err := gh.Deps.Core.DB.UpdateEmissionAndGNRE(ctx, database.UpdateEmissionAndGNREParams{
		ID:           pgtype.UUID{Bytes: p.EmissionId},
		Status:       pgtype.Text{String: "FINISHED", Valid: true},
		Message:      pgtype.Text{String: tecnospeedResponse.Sucess.Motivo, Valid: true},
		NumeroRecibo: tecnospeedResponse.Sucess.NumRecibo,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return fmt.Errorf("error updating GNRE", asynq.SkipRetry)
	}

	jsonGNRE, err := json.Marshal(updatedGNRE)
	if err != nil {
		return fmt.Errorf("error jsoing updatedGNRE")
	}

	_, err = t.ResultWriter().Write(jsonGNRE)
	if err != nil {
		return fmt.Errorf("failed to save result task", asynq.SkipRetry)
	}

	return nil
}

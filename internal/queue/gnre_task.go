package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/ent/gnreemission"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
	"github.com/lucasbonna/contafacil_api/internal/utils"
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
	ChaveNota     string
	XmlContent    string
	ClientDetails *schemas.ClientDetails
}

func (gh *GNREHandler) ProcessIssueGNRE(ctx context.Context, t *asynq.Task) error {
	var p IssueGNRETaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing GNRE: Client: %v", p.ClientDetails.Client.Name)

	tx, err := gh.Deps.Core.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	tecnospeedResponse, err := gh.Deps.External.TecnospeedService.IssueGNRE(p.XmlContent, "ContaFacil", p.ClientDetails.Client.Cnpj)
	if err != nil {
		if txErr := utils.FinishTask(tx, p.EmissionId, emission.StatusFAILED, "Erro ao emitir GNRE"); txErr != nil {
			return fmt.Errorf("failed to update emission: %v", txErr)
		}

		return fmt.Errorf("error: %v", err)
	}
	if tecnospeedResponse.Failure != nil {
		if tecnospeedResponse.Failure.Message != "" {
			if txErr := utils.FinishTask(tx, p.EmissionId, emission.StatusFAILED, tecnospeedResponse.Failure.Message); txErr != nil {
				return fmt.Errorf("failed to update emission: %v", txErr)
			}

			return fmt.Errorf("error: %v", err)
		}
	}

	fileBytes, err := gh.Deps.External.TecnospeedService.DownloadGNRE("ContaFacil", p.ClientDetails.Client.Cnpj, p.ChaveNota, tecnospeedResponse.Sucess.NumRecibo)
	if err != nil {
		log.Println("error gnre download", err)
		return fmt.Errorf("failed to downloadGNRE")
	}

	if len(fileBytes) == 0 {
		return fmt.Errorf("bad pdf download GNRE")
	}

	fileReader := bytes.NewReader(fileBytes)

	fileId := uuid.New()
	gh.Deps.Core.SM.Upload(fileReader, fileId)

	gnreUpdate, err := tx.GnreEmission.Update().
		Where(gnreemission.HasEmissionWith(emission.ID(p.EmissionId))).
		SetNumeroRecibo(tecnospeedResponse.Sucess.NumRecibo).
		SetPdf(fileId).
		Save(ctx)
	if err != nil || gnreUpdate == 0 {
		return fmt.Errorf("gnre emission update failed: %v (records affected: %d)", err, gnreUpdate)
	}

	fullEmission, err := tx.Emission.Query().
		Where(emission.ID(p.EmissionId)).
		WithGnreEmission().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to load emission relations: %v", err)
	}

	if fullEmission.Edges.GnreEmission == nil {
		return fmt.Errorf("critical error: gnre emission not found for emission %s", p.EmissionId)
	}

	if err := utils.FinishTask(tx, p.EmissionId, emission.StatusFINISHED, tecnospeedResponse.Sucess.Motivo); err != nil {
		return fmt.Errorf("emission update failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	ssePayload := SSEUpdatePayload{
		EmissionID: p.EmissionId,
		Status:     emission.StatusFINISHED,
		Message:    tecnospeedResponse.Sucess.Motivo,
		ClientID:   p.ClientDetails.Client.ID,
		UserID:     p.ClientDetails.User.ID,
	}

	task, err := NewTask(TypeSSEUpdate, ssePayload)
	if err != nil {
		log.Printf("failed to enqueue SSE update: %v", err)
	}

	_, err = gh.Deps.Core.AQ.Enqueue(task, asynq.Queue("IssueGNREQueue"), asynq.Retention(48*time.Hour))
	if err != nil {
		return err
	}

	return nil
}

package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/ent/gnreemission"
	"github.com/lucasbonna/contafacil_api/internal/app"
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
		_, txErr := tx.Emission.UpdateOneID(p.EmissionId).
			SetStatus(emission.StatusFAILED).
			SetMessage("Erro ao emitir GNRE").
			Save(ctx)
		if txErr != nil {
			return fmt.Errorf("failed to update emission: %v", txErr)
		}

		return fmt.Errorf("error: %v", err)
	}
	if tecnospeedResponse.Failure != nil {
		if tecnospeedResponse.Failure.Message != "" {
			_, txErr := tx.Emission.UpdateOneID(p.EmissionId).
				SetStatus(emission.StatusFAILED).
				SetMessage(tecnospeedResponse.Failure.Message).
				Save(ctx)
			if txErr != nil {
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

	fileReader := bytes.NewReader(fileBytes)

	gh.Deps.Core.SM.Upload(fileReader, uuid.New())

	_, err = tx.Emission.UpdateOneID(p.EmissionId).
		SetStatus(emission.StatusFINISHED).
		SetMessage(tecnospeedResponse.Sucess.Motivo).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("emission update failed: %v", err)
	}

	gnreUpdate, err := tx.GnreEmission.Update().
		Where(gnreemission.HasEmissionWith(emission.ID(p.EmissionId))).
		SetNumeroRecibo(tecnospeedResponse.Sucess.NumRecibo).
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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	gnreData := fullEmission.Edges.GnreEmission
	gnre := schemas.GNRE{
		ID:             fullEmission.ID,
		ClientId:       fullEmission.ClientID,
		UserId:         fullEmission.UserID,
		Message:        fullEmission.Message,
		EmissionType:   string(fullEmission.EmissionType),
		Status:         string(fullEmission.Status),
		ChaveNota:      gnreData.ChaveNota,
		CodBarrasGuia:  gnreData.CodBarrasGuia,
		ComprovantePDF: gnreData.ComprovantePdf,
		GuiaAmount:     gnreData.GuiaAmount,
		NumeroRecibo:   gnreData.NumeroRecibo,
		PDF:            gnreData.Pdf,
		XML:            gnreData.XML,
		Destinatario:   gnreData.Destinatario,
		NumNota:        gnreData.NumNota,
		CreatedAt:      fullEmission.CreatedAt,
		UpdatedAt:      fullEmission.UpdatedAt,
		DeletedAt:      fullEmission.DeletedAt,
	}

	jsonGNRE, err := json.Marshal(gnre)
	if err != nil {
		return fmt.Errorf("error marshaling gnre: %v", err)
	}

	if _, err = t.ResultWriter().Write(jsonGNRE); err != nil {
		return fmt.Errorf("failed to save result task: %v", err)
	}

	return nil
}

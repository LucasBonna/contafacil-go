package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/ent"
	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/queue"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

type EmissionHandlers struct {
	core *app.CoreDependencies
	ext  *app.ExternalDependencies
	int  *app.InternalDependencies
}

func NewEmissionHandlers(core *app.CoreDependencies, ext *app.ExternalDependencies, int *app.InternalDependencies) *EmissionHandlers {
	return &EmissionHandlers{
		core: core,
		ext:  ext,
		int:  int,
	}
}

func (eh *EmissionHandlers) HandlerListEmissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			clientID   uuid.UUID
			status     string
			startDate  time.Time
			endDate    time.Time
			parseError error
		)

		if idStr := c.Query("clientId"); idStr != "" {
			if clientID, parseError = uuid.Parse(idStr); parseError != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
				return
			}
		}

		status = c.Query("status")

		const dateLayout = "02-01-2006"
		if sdStr := c.Query("startDate"); sdStr != "" {
			if startDate, parseError = time.Parse(dateLayout, sdStr); parseError != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inicial inválido (DD-MM-AAAA)"})
				return
			}
		}

		if edStr := c.Query("endDate"); edStr != "" {
			if endDate, parseError = time.Parse(dateLayout, edStr); parseError != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data final inválido (DD-MM-AAAA)"})
				return
			}
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		if size < 1 || page < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros de paginação inválidos"})
			return
		}

		includeDeleted, _ := strconv.ParseBool(c.DefaultQuery("includeDeleted", "false"))

		baseQuery := eh.core.DB.Emission.Query()

		if clientID != uuid.Nil {
			baseQuery = baseQuery.Where(emission.ClientID(clientID))
		}

		if status != "" {
			baseQuery = baseQuery.Where(emission.StatusEQ(emission.Status(status)))
		}

		if !startDate.IsZero() {
			baseQuery = baseQuery.Where(emission.CreatedAtGTE(startDate))
		}

		if !endDate.IsZero() {
			baseQuery = baseQuery.Where(emission.CreatedAtLTE(endDate))
		}

		if !includeDeleted {
			baseQuery = baseQuery.Where(emission.DeletedAtIsNil())
		}

		// Get total count
		total, err := baseQuery.Count(c.Request.Context())
		if err != nil {
			log.Printf("Erro ao contar emissões: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao contar emissões"})
			return
		}

		// Apply pagination and include GnreEmission
		emissions, err := baseQuery.
			Offset(page * size).
			Limit(size).
			Order(ent.Desc(emission.FieldCreatedAt)).
			WithGnreEmission().
			All(c.Request.Context())
		if err != nil {
			log.Printf("Erro ao buscar emissões: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar emissões"})
			return
		}

		type GnreEmissionResponse struct {
			ID             uuid.UUID `json:"id"`
			XML            uuid.UUID `json:"xml"`
			Pdf            uuid.UUID `json:"pdf"`
			ComprovantePdf uuid.UUID `json:"comprovante_pdf"`
			GuiaAmount     float64   `json:"guia_amount"`
			NumeroRecibo   string    `json:"numero_recibo"`
			ChaveNota      string    `json:"chave_nota"`
			CodBarrasGuia  string    `json:"cod_barras_guia"`
			NumNota        string    `json:"num_nota"`
			Destinatario   string    `json:"destinatario"`
			CpfCnpj        string    `json:"cpf_cnpj"`
		}

		type EmissionResponse struct {
			ID           uuid.UUID             `json:"id"`
			EmissionType string                `json:"emission_type"`
			ClientID     uuid.UUID             `json:"client_id"`
			Message      string                `json:"message"`
			Status       string                `json:"status"`
			UserID       uuid.UUID             `json:"user_id"`
			CreatedAt    time.Time             `json:"created_at"`
			UpdatedAt    time.Time             `json:"updated_at"`
			DeletedAt    *time.Time            `json:"deleted_at"`
			GnreEmission *GnreEmissionResponse `json:"gnre_emission,omitempty"`
		}

		response := make([]EmissionResponse, len(emissions))
		for i, e := range emissions {
			res := EmissionResponse{
				ID:           e.ID,
				EmissionType: string(e.EmissionType),
				ClientID:     e.ClientID,
				Message:      e.Message,
				Status:       string(e.Status),
				UserID:       e.UserID,
				CreatedAt:    e.CreatedAt,
				UpdatedAt:    e.UpdatedAt,
			}
			if !e.DeletedAt.IsZero() {
				res.DeletedAt = &e.DeletedAt
			}
			if e.Edges.GnreEmission != nil {
				ge := e.Edges.GnreEmission
				res.GnreEmission = &GnreEmissionResponse{
					ID:             ge.ID,
					XML:            ge.XML,
					Pdf:            ge.Pdf,
					ComprovantePdf: ge.ComprovantePdf,
					GuiaAmount:     ge.GuiaAmount,
					NumeroRecibo:   ge.NumeroRecibo,
					ChaveNota:      ge.ChaveNota,
					CodBarrasGuia:  ge.CodBarrasGuia,
					NumNota:        ge.NumNota,
					Destinatario:   ge.Destinatario,
					CpfCnpj:        ge.CpfCnpj,
				}
			}
			response[i] = res
		}

		c.JSON(http.StatusOK, gin.H{
			"page":      page,
			"size":      size,
			"total":     total,
			"emissions": response,
		})
	}
}

func (eh *EmissionHandlers) IssueGNREHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientDetails := utils.GetClientDetails(c)

		mf, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}

		files := mf.File["files"]

		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
			return
		}

		// Keep track of all invoice numbers we have processed in this request
		processedNotas := make(map[string]bool)

		emissionIDs := make([]uuid.UUID, 0, len(files))

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to open file"})
				return
			}
			defer file.Close()

			xmlBytes, err := io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to read file"})
				return
			}

			validateResp, err := eh.int.XMLService.ValidateAndProcess(xmlBytes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// --- SKIP if we already processed this invoice number in this request
			if processedNotas[validateResp.NumNota] {
				log.Printf("Skipping file with repeated NumNota: %s", validateResp.NumNota)
				continue
			}
			processedNotas[validateResp.NumNota] = true

			// We must re-open the file if your storage manager needs a reader again.
			fileForUpload, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to reopen file, may be corrupted",
				})
				return
			}
			defer fileForUpload.Close()

			fileId := uuid.New()
			err = eh.core.SM.Upload(fileForUpload, fileId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			tx, err := eh.core.DB.Tx(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
				return
			}
			defer tx.Rollback()

			newEmission, err := tx.Emission.Create().
				SetID(uuid.New()).
				SetEmissionType("GNRE").
				SetClientID(clientDetails.Client.ID).
				SetUserID(clientDetails.User.ID).
				SetStatus(emission.StatusPROCESSING).
				SetMessage("Processando GNRE").
				Save(c.Request.Context())
			if err != nil {
				log.Println("error emission", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create emission"})
				return
			}

			_, err = tx.GnreEmission.Create().
				SetEmission(newEmission).
				SetID(newEmission.ID).
				SetXML(fileId).
				SetGuiaAmount(validateResp.IcmsValue).
				SetChaveNota(validateResp.ChaveNota).
				SetNumNota(validateResp.NumNota).
				SetDestinatario(validateResp.Destinatario).
				SetCpfCnpj(validateResp.CpfCnpj).
				Save(c.Request.Context())
			if err != nil {
				log.Println("error gnre emission", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create gnre emission"})
				return
			}

			if err := tx.Commit(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
				return
			}

			taskPayload := queue.IssueGNRETaskPayload{
				EmissionId:    newEmission.ID,
				XmlContent:    validateResp.ProcessedXML,
				ChaveNota:     validateResp.ChaveNota,
				ClientDetails: clientDetails,
			}

			task, err := queue.NewTask(string(queue.TypeIssueGNRE), taskPayload)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			info, err := eh.core.AQ.Enqueue(task, asynq.Queue("IssueGNREQueue"), asynq.Retention(48*time.Hour))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

			emissionIDs = append(emissionIDs, newEmission.ID)
		}

		c.JSON(http.StatusOK, gin.H{
			"total_enqueued": len(emissionIDs),
			"emission_ids":   emissionIDs,
		})
	}
}

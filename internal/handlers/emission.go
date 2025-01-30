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

		query := eh.core.DB.Emission.Query()

		if clientID != uuid.Nil {
			query = query.Where(emission.ClientID(clientID))
		}

		if status != "" {
			query = query.Where(emission.StatusEQ(emission.Status(status)))
		}

		if !startDate.IsZero() {
			query = query.Where(emission.CreatedAtGTE(startDate))
		}

		if !endDate.IsZero() {
			query = query.Where(emission.CreatedAtLTE(endDate))
		}

		if !includeDeleted {
			query = query.Where(emission.DeletedAtIsNil())
		}

		query = query.
			Offset(page * size).
			Limit(size).
			Order(ent.Desc(emission.FieldCreatedAt))

		emissions, err := query.All(c.Request.Context())
		if err != nil {
			log.Printf("Erro ao buscar emissões: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar emissões"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":      page,
			"size":      size,
			"emissions": emissions,
		})
	}
}

func (eh *EmissionHandlers) IssueGNREHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientDetails := utils.GetClientDetails(c)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "a file is required",
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to open file",
			})
			return
		}
		defer f.Close()

		xmlBytes, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable read file",
			})
			return
		}

		validateResp, err := eh.int.XMLService.ValidateAndProcess(xmlBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		fileId := uuid.New()
		err = eh.core.SM.Upload(f, fileId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		tx, err := eh.core.DB.Tx(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
			return
		}
		defer tx.Rollback()

		emission, err := tx.Emission.Create().
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
			SetEmission(emission).
			SetID(emission.ID).
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
			EmissionId:    emission.ID,
			XmlContent:    validateResp.ProcessedXML,
			ChaveNota:     validateResp.ChaveNota,
			ClientDetails: clientDetails,
		}

		task, err := queue.NewTask(string(queue.TypeIssueGNRE), taskPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		info, err := eh.core.AQ.Enqueue(task, asynq.Queue("IssueGNREQueue"), asynq.Retention(48*time.Hour))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

		c.JSON(http.StatusOK, emission)
	}
}

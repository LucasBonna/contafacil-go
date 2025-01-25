package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/database"
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
			clientId  pgtype.UUID
			status    pgtype.Text
			startDate pgtype.Timestamp
			endDate   pgtype.Timestamp
			err       error
		)

		if id := c.DefaultQuery("clientId", ""); id != "" {
			parsedUUID, err := uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid client id",
				})
				return
			}
			clientId = pgtype.UUID{Bytes: parsedUUID, Valid: true}
		} else {
			clientId = pgtype.UUID{Valid: false}
		}

		if s := c.DefaultQuery("status", ""); s != "" {
			status = pgtype.Text{String: s, Valid: true}
		} else {
			status = pgtype.Text{Valid: false}
		}

		const customDateLayout = "02-01-2006"
		if sd := c.DefaultQuery("startDate", ""); sd != "" {
			parsed, err := time.Parse(customDateLayout, sd)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format. Expected DD-MM-YYYY"})
				return
			}
			startDate = pgtype.Timestamp{Time: parsed, Valid: true}
		} else {
			startDate = pgtype.Timestamp{Valid: false}
		}

		if ed := c.DefaultQuery("endDate", ""); ed != "" {
			parsed, err := time.Parse(customDateLayout, ed)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format. Expected DD-MM-YYYY"})
				return
			}
			endDate = pgtype.Timestamp{Time: parsed, Valid: true}
		} else {
			endDate = pgtype.Timestamp{Valid: false}
		}

		includeDeleted := c.DefaultQuery("includeDeleted", "false")
		var includeDeletedBool bool
		if includeDeleted != "" {
			includeDeletedBool, err = strconv.ParseBool(includeDeleted)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid includeDeleted value"})
				return
			}
		}

		page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
		if err != nil || page < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page value"})
			return
		}

		size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
		if err != nil || size <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size value"})
			return
		}

		emissions, err := eh.core.DB.GetEmissionsByFilters(c.Request.Context(),
			database.GetEmissionsByFiltersParams{
				ClientID:       clientId,
				Status:         status,
				StartDate:      startDate,
				EndDate:        endDate,
				IncludeDeleted: includeDeletedBool,
				RowLimit:       int32(size),
				RowOffset:      int32(page * size),
			})
		if err != nil {
			log.Println("failed to fetch emissions: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch emissions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":    page,
			"size":    size,
			"content": emissions,
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
		log.Println("fileID", pgtype.UUID{Bytes: fileId, Valid: true})
		log.Println("icms", validateResp.IcmsValue)
		log.Println("chave", validateResp.ChaveNota)
		log.Println("numNota", validateResp.NumNota)
		log.Println("destinatario", validateResp.Destinatario)
		log.Println("userId", clientDetails.User.ID)
		log.Println("clientId", clientDetails.Client.ID)

		createdEmission, err := eh.core.DB.CreateGNREEmission(context.Background(), database.CreateGNREEmissionParams{
			XmlFileID:    pgtype.UUID{Bytes: fileId, Valid: true},
			EmissionType: pgtype.Text{String: "GNRE", Valid: true},
			GuiaAmount:   validateResp.IcmsValue,
			ChaveNota:    validateResp.ChaveNota,
			NumNota:      validateResp.NumNota,
			Destinatario: validateResp.Destinatario,
			ClientID:     clientDetails.Client.ID,
			Message:      pgtype.Text{String: "Processando GNRE", Valid: true},
			Status:       pgtype.Text{String: "PROCESSING", Valid: true},
			UserID:       clientDetails.User.ID,
		})
		if err != nil {
			// ignore error for now and get emission some other way
		}
		emission, err := eh.core.DB.GetGNREEmissionById(context.Background(), createdEmission.ID)

		parsedUUId, _ := uuid.FromBytes(emission.ID.Bytes[:])

		log.Println("parsedUUID", parsedUUId)

		taskPayload := queue.IssueGNRETaskPayload{
			EmissionId:    parsedUUId,
			XmlContent:    validateResp.ProcessedXML,
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

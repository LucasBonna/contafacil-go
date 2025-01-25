package middlewares

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/database"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func isJsonContentType(contentType string) bool {
	return contentType == "application/json"
}

func Logger(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		if strings.HasPrefix(c.Request.URL.Path, "/swagger") || c.Request.URL.Path == "/docs" {
			c.Next()
			return
		}

		t := time.Now()

		requestHeaders := formatHeaders(c.Request.Header)
		requestQuery := c.Request.URL.RawQuery
		requestParams := formatParams(c.Params)

		logParams := database.CreateAccessLogParams{
			ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Ip:             c.ClientIP(),
			Method:         c.Request.Method,
			Endpoint:       pgtype.Text{String: c.Request.URL.Path, Valid: true},
			RequestHeaders: pgtype.Text{String: requestHeaders, Valid: requestHeaders != ""},
			RequestParams:  pgtype.Text{String: requestParams, Valid: requestParams != ""},
			RequestQuery:   pgtype.Text{String: requestQuery, Valid: requestQuery != ""},
			CreatedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		contentType := c.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			fileUploadMsg := "Multipart form data (file upload)"
			logParams.RequestBody = pgtype.Text{String: fileUploadMsg, Valid: fileUploadMsg != ""}
		} else if !isJsonContentType(contentType) {
			nonJSONMsg := fmt.Sprintf("Non-JSON content type: %s", contentType)
			logParams.RequestBody = pgtype.Text{String: nonJSONMsg, Valid: nonJSONMsg != ""}
			logParams.RequestHeaders = pgtype.Text{String: nonJSONMsg, Valid: nonJSONMsg != ""}
			logParams.ResponseHeaders = pgtype.Text{String: nonJSONMsg, Valid: nonJSONMsg != ""}
			logParams.ResponseBody = pgtype.Text{String: nonJSONMsg, Valid: nonJSONMsg != ""}
		} else {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		insertedLog, err := deps.Core.DB.CreateAccessLog(ctx, logParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		responseTime := time.Since(t).String()
		responseBody := blw.body.String()
		responseHeaders := formatHeaders(c.Writer.Header())

		updatedLogParams := database.UpdateAccessLogParams{
			ID:              insertedLog.ID,
			Ip:              insertedLog.Ip,
			Method:          insertedLog.Method,
			Endpoint:        insertedLog.Endpoint,
			RequestBody:     insertedLog.RequestBody,
			RequestHeaders:  insertedLog.RequestHeaders,
			RequestQuery:    insertedLog.RequestQuery,
			RequestParams:   insertedLog.RequestParams,
			ResponseBody:    pgtype.Text{String: responseBody, Valid: true},
			ResponseHeaders: pgtype.Text{String: responseHeaders, Valid: true},
			ResponseTime:    pgtype.Text{String: responseTime, Valid: responseTime != ""},
			StatusCode:      pgtype.Int4{Int32: int32(c.Writer.Status()), Valid: true},
			UpdatedAt:       pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		_, err = deps.Core.DB.UpdateAccessLog(ctx, updatedLogParams)
		if err != nil {
			log.Printf("error updating accesslog: %v", err)
		}
	}
}

func formatParams(params gin.Params) string {
	var buffer bytes.Buffer
	for _, p := range params {
		buffer.WriteString(p.Key)
		buffer.WriteString("=")
		buffer.WriteString(p.Value)
		buffer.WriteString("&")
	}
	if buffer.Len() > 0 {
		buffer.Truncate(buffer.Len() - 1)
	}
	return buffer.String()
}

func formatHeaders(headers map[string][]string) string {
	var buffer bytes.Buffer
	for key, values := range headers {
		buffer.WriteString(key)
		buffer.WriteString(": ")
		buffer.WriteString(strings.Join(values, ", "))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

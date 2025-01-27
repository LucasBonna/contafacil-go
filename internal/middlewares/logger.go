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

	"github.com/lucasbonna/contafacil_api/internal/app"
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

		logParams := deps.Core.DB.AccessLog.Create().
			SetID(uuid.New()).
			SetIP(c.ClientIP()).
			SetMethod(c.Request.Method).
			SetEndpoint(c.Request.URL.Path).
			SetRequestHeaders(formatHeaders(c.Request.Header)).
			SetRequestParams(formatParams(c.Params)).
			SetRequestQuery(c.Request.URL.RawQuery).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now())

		contentType := c.GetHeader("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			logParams.SetRequestBody("Multipart form data (file upload)")
		} else if !isJsonContentType(contentType) {
			nonJSONMsg := fmt.Sprintf("Non-JSON content type: %s", contentType)
			logParams.SetRequestBody(nonJSONMsg)
		} else {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			logParams.SetRequestBody(string(bodyBytes))
		}

		insertedLog, err := logParams.Save(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		responseTime := time.Since(t).String()

		update := deps.Core.DB.AccessLog.UpdateOneID(insertedLog.ID).
			SetResponseHeaders(formatHeaders(c.Writer.Header())).
			SetResponseBody(blw.body.String()).
			SetResponseTime(responseTime).
			SetStatusCode(c.Writer.Status()).
			SetUpdatedAt(time.Now())

		if contentType != "" && !isJsonContentType(contentType) {
			update.SetResponseHeaders("Non-JSON content type: " + contentType)
		}

		if _, err := update.Save(ctx); err != nil {
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

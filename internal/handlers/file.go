package handlers

import (
	"archive/zip"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

const MaxFileSize = 10 << 20

func HandlerUploadFile(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

		if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "files too large or invalid"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid file",
			})
			return
		}

		fileId := uuid.New()
		fileName := file.Filename
		contentType := file.Header.Get("Content-Type")

		openedFile, err := file.Open()
		if err != nil {
			log.Printf("error opening file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading file"})
			return
		}
		defer openedFile.Close()

		// Upload to storage
		if err := deps.Core.SM.Upload(openedFile, fileId); err != nil {
			log.Printf("error uploading file to storage: %v", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error uploading file to storage"},
			)
			return
		}

		createdFile, err := deps.Core.DB.CreateFile(c.Request.Context(), database.CreateFileParams{
			ID:          pgtype.UUID{Bytes: fileId, Valid: true},
			FileName:    fileName,
			Extension:   getFileExtension(fileName),
			ContentType: contentType,
			FilePath:    fileId.String(),
			CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		})
		if err != nil {
			log.Printf("error saving file meteadata on db: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving file on db"})
			return
		}

		c.JSON(http.StatusOK, createdFile)
	}
}

func HandlerDownloadFile(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileId := c.Param("fileId")
		if fileId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing file Id",
			})
			return
		}

		parsedFileId, err := uuid.Parse(fileId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid uuid",
			})
		}

		fileBytes, err := deps.Core.SM.Download(parsedFileId)
		if err != nil {
			log.Printf("error downloading file from storage: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error downloading file from storage",
			})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileId+".pdf"))
		c.Header("Content-Type", "application/octet-stream")

		c.Data(http.StatusOK, "application/octet-stream", fileBytes)
	}
}

func HandlerDownloadBatch(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body schemas.DownloadBatchFileSchema
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if len(body.FileIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no valid file ids"})
			return
		}

		tempFile, err := os.CreateTemp("", "batch_download_*.zip")
		if err != nil {
			log.Printf("error creating temp zip file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		zipWriter := zip.NewWriter(tempFile)
		defer zipWriter.Close()

		for _, fileId := range body.FileIds {
			fileBytes, err := deps.Core.SM.Download(fileId)
			if err != nil {
				log.Printf("error downloading file from storage: %v", err)
				zipWriter.Close()
				os.Remove(tempFile.Name())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("failed to download file with ID: %s", fileId.String()),
				})
				return
			}

			fileName := fmt.Sprintf("%s.pdf", fileId.String())
			zipFile, err := zipWriter.Create(fileName)
			if err != nil {
				log.Printf("error adding file to zip: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating zip file"})
				return
			}

			if _, err := zipFile.Write(fileBytes); err != nil {
				log.Printf("error writing to zip file: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing to zip file"})
				return
			}
		}

		zipWriter.Close()

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", `attachment; filename="batch_download.zip"`)
		c.File(tempFile.Name())
	}
}

func getFileExtension(fileName string) string {
	if idx := strings.LastIndex(fileName, "."); idx != -1 {
		return fileName[idx+1:]
	}

	return ""
}

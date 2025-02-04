package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"

	"github.com/lucasbonna/contafacil_api/ent"
	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func GetClientDetails(c *gin.Context) *schemas.ClientDetails {
	clientDetailsStored, exists := c.Get("clientDetails")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	clientDetails, ok := clientDetailsStored.(*schemas.ClientDetails)
	if !ok || clientDetails == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	return clientDetails
}

func FinishTask(tx *ent.Tx, emissionId uuid.UUID, status emission.Status, message string) error {
	_, err := tx.Emission.UpdateOneID(emissionId).
		SetStatus(status).
		SetMessage(message).
		Save(context.Background())
	return err
}

func GenerateAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16
	key := make([]byte, length)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}

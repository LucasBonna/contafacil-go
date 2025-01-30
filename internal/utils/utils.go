package utils

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/ent"
	"github.com/lucasbonna/contafacil_api/ent/emission"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

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

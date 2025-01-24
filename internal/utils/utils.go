package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

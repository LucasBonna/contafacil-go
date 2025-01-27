package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/ent/clients"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type ClientHandlers struct {
	Core *app.CoreDependencies
}

func NewClientHandlers(core *app.CoreDependencies) *ClientHandlers {
	return &ClientHandlers{
		Core: core,
	}
}

func (rh ClientHandlers) CreateClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input schemas.CreateClientInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := rh.Core.DB.Clients.Create().
			SetName(input.Name).
			SetCnpj(input.Cnpj).
			SetRole(input.Role).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create client"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":   client.ID,
			"name": client.Name,
			"cnpj": client.Cnpj,
			"role": client.Role,
		})
	}
}

func (rh ClientHandlers) GetClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
			return
		}

		client, err := rh.Core.DB.Clients.Query().
			Where(clients.ID(id), clients.DeletedAtIsNil()).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   client.ID,
			"name": client.Name,
			"cnpj": client.Cnpj,
			"role": client.Role,
		})
	}
}

func (rh ClientHandlers) ListAllClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", "10")

		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}

		offset := (pageInt - 1) * pageSizeInt

		query := rh.Core.DB.Clients.Query().
			Where(clients.DeletedAtIsNil())

		total, err := query.Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count clients"})
			return
		}

		clients, err := query.
			Offset(offset).
			Limit(pageSizeInt).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clients"})
			return
		}

		var response []gin.H
		for _, client := range clients {
			response = append(response, gin.H{
				"id":   client.ID,
				"name": client.Name,
				"cnpj": client.Cnpj,
				"role": client.Role,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
			"pagination": gin.H{
				"page":       pageInt,
				"pageSize":   pageSizeInt,
				"totalItems": total,
				"totalPages": int(math.Ceil(float64(total) / float64(pageSizeInt))),
			},
		})
	}
}

func (rh ClientHandlers) UpdateClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
			return
		}

		var input schemas.UpdateClientInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := rh.Core.DB.Clients.UpdateOneID(id)
		if input.Name != nil {
			update.SetName(*input.Name)
		}
		if input.Cnpj != nil {
			update.SetCnpj(*input.Cnpj)
		}
		if input.Role != nil {
			update.SetRole(*input.Role)
		}

		client, err := update.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update client"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   client.ID,
			"name": client.Name,
			"cnpj": client.Cnpj,
			"role": client.Role,
		})
	}
}

func (rh ClientHandlers) DeleteClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
			return
		}

		_, err = rh.Core.DB.Clients.UpdateOneID(id).
			SetDeletedAt(time.Now()).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete client"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

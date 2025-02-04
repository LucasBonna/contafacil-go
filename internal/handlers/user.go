package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/lucasbonna/contafacil_api/ent/user"
	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/schemas"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

type UserHandlers struct {
	Core *app.CoreDependencies
}

func NewUserHandlers(core *app.CoreDependencies) *UserHandlers {
	return &UserHandlers{
		Core: core,
	}
}

func (rh UserHandlers) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input schemas.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		user, err := rh.Core.DB.User.Create().
			SetID(uuid.New()).
			SetUsername(input.Username).
			SetPassword(string(hashedPassword)).
			SetAPIKey(utils.GenerateAPIKey()).
			SetRole(input.Role).
			SetClientID(input.ClientID).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"role":      user.Role,
			"clientId":  user.ClientID,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"deletedAt": user.DeletedAt,
		})
	}
}

func (rh UserHandlers) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientDetails := utils.GetClientDetails(c)
		if clientDetails == nil {
			return
		}

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		if clientDetails.User.ID != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		user, err := rh.Core.DB.User.Query().
			Where(user.ID(id), user.DeletedAtIsNil()).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"role":      user.Role,
			"clientId":  user.ClientID,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"deletedAt": user.DeletedAt,
		})
	}
}

func (rh UserHandlers) ListAllUsers() gin.HandlerFunc {
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

		clientIdStr := c.Query("clientId")
		query := rh.Core.DB.User.Query().Where(user.DeletedAtIsNil())

		if clientIdStr != "" {
			clientId, err := uuid.Parse(clientIdStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
				return
			}
			query = query.Where(user.ClientID(clientId))
		}

		total, err := query.Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
			return
		}

		users, err := query.
			Offset(offset).
			Limit(pageSizeInt).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
			return
		}

		var response []gin.H
		for _, user := range users {
			response = append(response, gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"role":      user.Role,
				"clientId":  user.ClientID,
				"createdAt": user.CreatedAt,
				"updatedAt": user.UpdatedAt,
				"deletedAt": user.DeletedAt,
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

func (rh UserHandlers) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		clientDetails := utils.GetClientDetails(c)

		if clientDetails.User.ID != id && clientDetails.User.Role != string(user.RoleADMIN) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you are not authorized to update this user",
			})
			return
		}

		var input schemas.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if clientDetails.User.ID != id && (input.Password != nil || input.Username != nil) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you cant update another user credentials"})
			return
		}

		update := rh.Core.DB.User.UpdateOneID(id)
		if input.Username != nil {
			update.SetUsername(*input.Username)
		}
		if input.Password != nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}
			update.SetPassword(string(hashedPassword))
		}
		if input.Role != nil {
			update.SetRole(*input.Role)
		}

		user, err := update.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"role":      user.Role,
			"clientId":  user.ClientID,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"deletedAt": user.DeletedAt,
		})
	}
}

func (rh UserHandlers) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		_, err = rh.Core.DB.User.UpdateOneID(id).
			SetDeletedAt(time.Now()).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

package schemas

import (
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/ent/user"
)

type CreateUserInput struct {
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Role     user.Role `json:"role" binding:"required"`
	ClientID uuid.UUID `json:"clientId" binding:"required"`
}

type UpdateUserInput struct {
	Username *string    `json:"username,omitempty"`
	Password *string    `json:"password,omitempty"`
	Role     *user.Role `json:"role,omitempty"`
}

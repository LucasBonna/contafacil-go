package schemas

import (
	"github.com/lucasbonna/contafacil_api/ent/clients"
)

type CreateClientInput struct {
	Name string       `json:"name" binding:"required"`
	Cnpj string       `json:"cnpj" binding:"required"`
	Role clients.Role `json:"role" binding:"required"`
}

type UpdateClientInput struct {
	Name *string       `json:"name,omitempty"`
	Cnpj *string       `json:"cnpj,omitempty"`
	Role *clients.Role `json:"role,omitempty"`
}

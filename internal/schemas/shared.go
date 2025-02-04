package schemas

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	ApiKey    string
	Role      string
	ClientID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Client struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cnpj      string    `json:"cnpj"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ClientDetails struct {
	User   User
	Client Client
}

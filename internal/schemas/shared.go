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
	ID        uuid.UUID
	Name      string
	Cnpj      string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ClientDetails struct {
	User   User
	Client Client
}

package schemas

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID
	Username  string
	ApiKey    string
	Role      pgtype.Text
	ClientID  pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type Client struct {
	ID        pgtype.UUID
	Name      string
	Cnpj      string
	Role      pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type ClientDetails struct {
	User   User
	Client Client
}

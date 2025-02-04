package schemas

import "github.com/google/uuid"

type JWTToken struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	ApiKey   string    `json:"apiKey"`
	Role     string    `json:"role"`
	Client   Client    `json:"client"`
}

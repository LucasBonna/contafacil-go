package schemas

import (
	"github.com/google/uuid"
	"github.com/lucasbonna/contafacil_api/ent/emission"
)

type SSEMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type EmissionUpdate struct {
	EmissionID uuid.UUID       `json:"emission_id"`
	Status     emission.Status `json:"status"`
	Message    string          `json:"message"`
}

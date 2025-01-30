package schemas

import (
	"time"

	"github.com/google/uuid"
)

type GNRE struct {
	ID             uuid.UUID
	ClientId       uuid.UUID
	UserId         uuid.UUID
	Message        string
	EmissionType   string
	Status         string
	ChaveNota      string
	CodBarrasGuia  string
	ComprovantePDF uuid.UUID
	GuiaAmount     float64
	NumeroRecibo   string
	PDF            uuid.UUID
	XML            uuid.UUID
	Destinatario   string
	NumNota        string
	CpfCnpj        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

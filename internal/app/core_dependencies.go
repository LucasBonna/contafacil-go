package app

import (
	"github.com/go-resty/resty/v2"

	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
	"github.com/lucasbonna/contafacil_api/internal/storage"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

type CoreDependencies struct {
	DB     *database.Queries
	Rabbit *rabbitmq.RabbitMQ
	QH     *utils.QueueHelper
	SM     storage.StorageManager
	RC     *resty.Client
}

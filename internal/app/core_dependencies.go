package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/hibiken/asynq"

	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
	"github.com/lucasbonna/contafacil_api/internal/storage"
)

type CoreDependencies struct {
	DB     *database.Queries
	AQ     *asynq.Client
	Rabbit *rabbitmq.RabbitMQ
	SM     storage.StorageManager
	RC     *resty.Client
}

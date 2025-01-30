package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/hibiken/asynq"
	"github.com/lucasbonna/contafacil_api/ent"

	"github.com/lucasbonna/contafacil_api/internal/storage"
)

type CoreDependencies struct {
	DB *ent.Client
	AQ *asynq.Client
	SM storage.StorageManager
	RC *resty.Client
}

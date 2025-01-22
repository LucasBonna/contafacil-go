package main

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
	"github.com/lucasbonna/contafacil_api/internal/server"
	"github.com/lucasbonna/contafacil_api/internal/services"
	"github.com/lucasbonna/contafacil_api/internal/storage"
	"github.com/lucasbonna/contafacil_api/internal/storage/r2"
	"github.com/lucasbonna/contafacil_api/internal/utils"
	"github.com/lucasbonna/contafacil_api/internal/worker"
)

func main() {
	config.LoadEnvs()

	rabbit, err := rabbitmq.NewRabbitMQ(config.RabbitMQUrl)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}

	queueHelper := utils.NewQueueHelper(rabbit)

	r2Client, err := r2.NewR2Client(
		config.StorageAccessKeyId,
		config.StorageAccessKeySecret,
		config.StorageRegion,
		config.StorageAccountId,
	)
	if err != nil {
		log.Fatalf("error connecting to r2: %v", err)
	}

	storageManager := storage.SetStorage(r2Client)

	// Create db connection
	dbConn, err := database.ConnectToDB(config.Db_url)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Get queries
	queries := database.New(dbConn)

	// Create Resty Client
	restyClient := resty.New()
	restyClient.SetTimeout(60 * time.Second)

	core_deps := &app.CoreDependencies{
		DB:     queries,
		Rabbit: rabbit,
		QH:     queueHelper,
		SM:     storageManager,
		RC:     restyClient,
	}

	tecnospeedService := services.NewTecnospeedService(restyClient, config.TSUsername, config.TSPassword, config.TSBaseUrl)

	external_deps := &app.ExternalDependencies{
		TecnospeedService: tecnospeedService,
	}

	xmlService := services.NewXmlService()

	internal_deps := &app.InternalDependencies{
		XMLService: xmlService,
	}

	deps := &app.Dependencies{
		Core:     *core_deps,
		External: *external_deps,
		Internal: *internal_deps,
	}

	// Criar dispatcher e iniciar handlers
	dispatcher := worker.NewDispatcher()
	dispatcher.RegisterHandler("IssueGNRE", &worker.IssueGNREHandler{})
	log.Println("dispatchers registrados")

	// Iniciar workers
	go worker.StartWorkers(rabbit, dispatcher)

	// Inicar worker de retries
	go worker.StartRetryWorker(rabbit)

	// Iniciar servidor HTTP
	server := server.NewServer(config.Db_url, rabbit, deps)
	server.StartServer()
}

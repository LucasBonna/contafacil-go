package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hibiken/asynq"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/queue"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
	"github.com/lucasbonna/contafacil_api/internal/server"
	"github.com/lucasbonna/contafacil_api/internal/services"
	"github.com/lucasbonna/contafacil_api/internal/storage"
	"github.com/lucasbonna/contafacil_api/internal/storage/r2"
)

func main() {
	config.InitEnvs()

	log.Println("teste", config.Env.RabbitMQUrl)

	rabbit, err := rabbitmq.NewRabbitMQ(config.Env.RabbitMQUrl)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}

	r2Client, err := r2.NewR2Client(
		config.Env.StorageAccessKeyId,
		config.Env.StorageAccessKeySecret,
		config.Env.StorageRegion,
		config.Env.StorageAccountId,
	)
	if err != nil {
		log.Fatalf("error connecting to r2: %v", err)
	}

	storageManager := storage.SetStorage(r2Client)

	// Create db connection
	dbConn, err := database.ConnectToDB(config.Env.Db_url)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Get queries
	queries := database.New(dbConn)

	// Create Resty Client
	restyClient := resty.New()
	restyClient.SetTimeout(60 * time.Second)

	// Create Asynq Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env.RedisAddr})
	defer asynqClient.Close()

	core_deps := &app.CoreDependencies{
		DB:     queries,
		AQ:     asynqClient,
		Rabbit: rabbit,
		SM:     storageManager,
		RC:     restyClient,
	}

	tecnospeedService := services.NewTecnospeedService(restyClient, config.Env.TSUsername, config.Env.TSPassword, config.Env.TSBaseUrl)

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

	// Configurar o servidor Asynq
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Env.RedisAddr},
		asynq.Config{
			Concurrency: 10, // Número de workers
			Queues: map[string]int{
				"default":  1,
				"critical": 2,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeIssueGNRE, queue.HandleIssueGNRETask)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run Asynq server: %v", err)
		}
	}()
	log.Println("Asynq server initialized")

	go func() {
		server := server.NewServer(config.Env.Db_url, rabbit, deps)
		server.StartServer()
	}()
	log.Println("Gin server iniciado")

	// Esperar por sinais de interrupção para shutdown graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Recebido sinal %s, iniciando shutdown", sig)

	// Shutdown do Asynq
	srv.Shutdown()

	log.Println("Shutdown completo")
}

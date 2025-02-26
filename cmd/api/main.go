package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hibiken/asynq"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"

	"github.com/lucasbonna/contafacil_api/internal/app"
	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/queue"
	"github.com/lucasbonna/contafacil_api/internal/server"
	"github.com/lucasbonna/contafacil_api/internal/services"
	"github.com/lucasbonna/contafacil_api/internal/sse"
	"github.com/lucasbonna/contafacil_api/internal/storage"
	"github.com/lucasbonna/contafacil_api/internal/storage/r2"
)

func main() {
	config.InitEnvs()

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
	dbConn := database.ConnectToDB(config.Env.DB_ConnStr)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer dbConn.Close()

	// Create Resty Client
	restyClient := resty.New()
	restyClient.SetTimeout(60 * time.Second)

	// Create Asynq Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env.RedisAddr})
	defer asynqClient.Close()

	// Create Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Env.RedisAddr,
		Password: "", // se necessário
		DB:       0,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	defer redisClient.Close()

	// Create SSEManager
	sseManager := sse.NewManager(redisClient)

	coreDeps := &app.CoreDependencies{
		DB:     dbConn,
		AQ:     asynqClient,
		SM:     storageManager,
		RC:     restyClient,
		SSEMgr: sseManager,
	}

	tecnospeedService := services.NewTecnospeedService(restyClient, config.Env.TSUsername, config.Env.TSPassword, config.Env.TSBaseUrl)

	externalDeps := &app.ExternalDependencies{
		TecnospeedService: tecnospeedService,
	}

	xmlService := services.NewXmlService()

	internalDeps := &app.InternalDependencies{
		XMLService: xmlService,
	}

	deps := &app.Dependencies{
		Core:     *coreDeps,
		External: *externalDeps,
		Internal: *internalDeps,
	}

	if config.Env.Type == "worker" {
		// Configuração do Asynq Server
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: "redis:6379"},
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"IssueGNREQueue":    5,
					"SSEEmissionUpdate": 10,
					"critical":          2,
				},
			},
		)

		gnreHandler := queue.NewGNREHandler(deps)
		sseHandler := queue.NewSSEHandler(deps)

		mux := asynq.NewServeMux()
		mux.HandleFunc(queue.TypeIssueGNRE, gnreHandler.ProcessIssueGNRE)
		mux.HandleFunc(queue.TypeSSEEmissionUpdate, sseHandler.ProcessSSEUpdate)

		// Canal para capturar erros do servidor Asynq
		errChan := make(chan error, 1)
		go func() {
			if err := srv.Run(mux); err != nil {
				errChan <- err
			}
		}()

		log.Println("Asynq server initialized")

		// Esperar por sinais OU erro do servidor
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			log.Printf("Recebido sinal %s, iniciando shutdown", sig)
		case err := <-errChan:
			log.Printf("Erro no servidor Asynq: %v", err)
		}

		srv.Shutdown()
		log.Println("Shutdown completo do worker")
		return // Sair após shutdown

	}
	server := server.NewServer(deps)

	// Canal para servidor Gin
	errChan := make(chan error, 1)
	go func() {
		server.StartServer()
	}()

	log.Println("Gin server iniciado")

	// Esperar por sinais OU erro do servidor
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Recebido sinal %s, iniciando shutdown", sig)
	case err := <-errChan:
		log.Printf("Erro no servidor Gin: %v", err)
	}

	log.Println("Shutdown completo do servidor")
}

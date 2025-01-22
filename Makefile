include .env

dsn := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
migrationPath := ./migrations

goose-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) up

goose-create:
	@if [ -z "$(name)" ]; then echo "Error: Please specify a name for the migration. Usage: make goose-create name=migration_name"; exit 1; fi
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) create $(name) sql

all-up:
	@docker-compose up -d
	sleep 5
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) up
	@sqlc generate

down:
	@docker-compose down

rs:
	@make down
	@make all-up

include .env

all:
	@go generate ./ent
	@docker-compose up -d

down:
	@docker-compose down

rs:
	@make down
	@make all

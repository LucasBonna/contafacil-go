package database

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	"github.com/lucasbonna/contafacil_api/ent"
)

func ConnectToDB(connectionString string) *ent.Client {
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

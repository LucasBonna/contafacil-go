package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucasbonna/contafacil_api/ent"
)

func ConnectToDB(connectionString string) *ent.Client {
	// Instead of "postgres", use "pgx" here:
	client, err := ent.Open("pgx", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to pgx: %v", err)
	}

	// Now create schema
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

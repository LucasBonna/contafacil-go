package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/lucasbonna/contafacil_api/ent"
)

func ConnectToDB(host string, port string, user string, dbname string, password string) *ent.Client {
	log.Println(host)
	log.Println(port)
	log.Println(user)
	log.Println(dbname)
	log.Println(password)
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

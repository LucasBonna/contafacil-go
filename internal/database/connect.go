package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB(dsn string) (*pgx.Conn, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
  defer cancel()

  conn, err := pgx.Connect(ctx, dsn)
  if err != nil {
    return nil, err
  }
  return conn, nil
}

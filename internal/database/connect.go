package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog/log"
)

func ConnectToDB(dsn string) (*pgx.Conn, error) {
	// Parsear configuração
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear DSN: %w", err)
	}

	// Configurar logger customizado
	config.Tracer = &tracelog.TraceLog{
		Logger:   &QueryLogger{},
		LogLevel: tracelog.LogLevelDebug, // Loga todas as queries
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Conectar com configuração
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	return conn, nil
}

type QueryLogger struct{}

func (ql *QueryLogger) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	if msg == "Query" {
		query := data["sql"].(string)
		args := data["args"].([]any)
		duration := data["time"].(time.Duration)

		log.Debug().
			Str("query", query).
			Interface("args", args).
			Dur("duration", duration).
			Msg("SQL Executado")
	}
}

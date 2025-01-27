package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/floroz/go-social/internal/env"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func ConnectDb() (*sql.DB, error) {
	postgresConnection := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		env.GetEnvValue("DB_USER"),
		env.GetEnvValue("DB_PASSWORD"),
		env.GetEnvValue("DB_HOST"),
		env.GetEnvValue("DB_NAME"),
	)

	db, err := sql.Open("postgres", postgresConnection)
	if err != nil {
		return nil, err
	}

	const connectionTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(60 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	log.Info().Msg("database connection pool established")
	return db, nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
	_ "github.com/lib/pq"
)

func connectDb() (*sql.DB, error) {
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

	log.Println("database connection pool established")
	return db, nil
}

func main() {
	env.MustLoadEnv(".env.local")

	db, err := connectDb()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		panic(err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	config := &config{
		port: env.GetEnvValue("PORT"),
	}

	app := &application{
		config:      config,
		userService: userService,
	}

	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

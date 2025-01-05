package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
	_ "github.com/lib/pq"
)

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", env.GetEnvValue("DB_USER"), env.GetEnvValue("DB_PASSWORD"), env.GetEnvValue("DB_HOST"), env.GetEnvValue("DB_NAME")))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	env.MustLoadEnv(".env.local")

	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	config := &config{
		port: env.GetEnvValue("PORT"),
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	app := &application{
		config,
		userService,
	}

	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

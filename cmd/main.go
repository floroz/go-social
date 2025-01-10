package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/floroz/go-social/cmd/api"
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

	log.Info().Msg("database connection pool established")
	return db, nil
}

func main() {
	env.MustLoadEnv(".env.local")

	// crash immediately if there's no JWT_SECRET
	if env.GetJWTSecret() == "" {
		panic("fatal: JWT_SECRET is required but not set in env.")
	}

	db, err := connectDb()
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)

	config := &api.Config{
		Port: env.GetEnvValue("PORT"),
	}

	app := &api.Application{
		Config:         config,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.Config.Port),
		Handler:      app.Routes(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	log.Info().Msgf("Starting server on %s", app.Config.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("server error")
	}
}

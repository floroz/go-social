package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/floroz/go-social/cmd/api"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func runMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	return nil
}

func startAPIServer(db *sql.DB) func() {
	env.MustLoadEnv("../../.env.local")

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo, commentRepo)

	authService := services.NewAuthService(userRepo)

	config := &api.Config{
		Port: env.GetEnvValue("PORT"),
	}

	app := &api.Application{
		Config:         config,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		AuthService:    authService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.Config.Port),
		Handler:      app.Routes(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	log.Info().Msgf("Started server on %s", app.Config.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server error")
		}
	}()

	time.Sleep(2 * time.Second)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
		log.Info().Msg("Server gracefully stopped")
	}
}

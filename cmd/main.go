package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/floroz/go-social/cmd/api"
	"github.com/floroz/go-social/cmd/database"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
)

func main() {
	env.MustLoadEnv(".env.local")

	// crash immediately if there's no JWT_SECRET
	if env.GetJWTSecret() == "" {
		panic("fatal: JWT_SECRET is required but not set in env.")
	}

	db, err := database.ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	defer db.Close()

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

	log.Info().Msgf("Starting server on %s", app.Config.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("server error")
	}
}

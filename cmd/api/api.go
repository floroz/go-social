package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config *config
}

type config struct {
	port string
}

func (app *application) run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.port),
		Handler:      app.routes(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("Starting server on %s\n", app.config.port)
	return server.ListenAndServe()
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
		})
	})

	return r
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user created"))
}

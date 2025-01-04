package main

import "net/http"

type application struct {
	config *config
}

type config struct {
	address string
}

func (app *application) run() error {
	server := &http.Server{
		Addr:    app.config.address,
		Handler: app.routes(),
	}

	return server.ListenAndServe()
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// posts

	// users

	// auth

	return mux
}

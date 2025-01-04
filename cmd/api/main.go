package main

import (
	"log"

	"github.com/floroz/go-social/internal/env"
)

func main() {
	env.MustLoadEnv()

	config := &config{
		port: env.GetEnvValue("ADDRESS", ":8080"),
	}

	app := &application{
		config: config,
	}

	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

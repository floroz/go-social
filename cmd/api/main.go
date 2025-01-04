package main

import "log"

func main() {
	config := &config{
		address: ":8080", // TODO: load from env file
	}

	app := &application{
		config: config,
	}

	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}

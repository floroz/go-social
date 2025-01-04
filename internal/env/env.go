package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvValue(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func MustLoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		panic(err)
	}
}

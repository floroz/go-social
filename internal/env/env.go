package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Key %s not found in ENV", key)
		return ""
	}
	return value
}

func MustLoadEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatalf("Error loading .env file")
		panic(err)
	}
}

package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func GetEnvValue(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warn().Msgf("Key %s not found in ENV", key)
		return ""
	}
	return value
}

func MustLoadEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
		panic(err)
	}
}

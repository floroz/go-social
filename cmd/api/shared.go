package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/floroz/go-social/internal/domain"
	"github.com/rs/zerolog/log"
)

type errorResponse struct {
	Errors []domain.ErrorDetail `json:"errors"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

func readJSON(source io.Reader, dest any) error {
	maxBytes := int64(1 << 20) // 1MB
	// Prevent reading too much data from the request body
	// to avoid potential denial of service attacks
	source = io.LimitReader(source, maxBytes)
	decoder := json.NewDecoder(source)

	// Ensure that the request body does not contain unknown fields
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	errors := []domain.ErrorDetail{
		{Message: message},
	}
	errorResponse := errorResponse{
		Errors: errors,
	}
	writeJSON(w, status, errorResponse)
}

func handleErrors(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("error")
	switch e := err.(type) {
	case *domain.ValidationError:
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case *domain.InternalServerError:
		writeJSONError(w, e.StatusCode, err.Error())
	case *domain.BadRequestError:
		writeJSONError(w, e.StatusCode, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

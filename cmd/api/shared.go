package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/floroz/go-social/internal/domain"
)

type errorResponse struct {
	Errors []domain.ErrorDetail `json:"errors"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
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
	switch err.(type) {
	case *domain.ValidationError:
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case *domain.InternalServerError:
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	case *domain.BadRequestError:
		writeJSONError(w, http.StatusBadRequest, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

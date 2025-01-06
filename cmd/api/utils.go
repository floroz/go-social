package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/floroz/go-social/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(source io.Reader, dest any) error {
	return json.NewDecoder(source).Decode(dest)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	data := map[string]string{"error": message}
	writeJSON(w, status, data)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *domain.ValidationError:
		errorResponse(w, http.StatusBadRequest, e.Error())
	case *domain.ConflictError:
		errorResponse(w, http.StatusConflict, e.Error())
	case *domain.InternalServerError:
		errorResponse(w, http.StatusInternalServerError, e.Error())
	case *domain.BadRequestError:
		errorResponse(w, http.StatusBadRequest, e.Error())
	default:
		errorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}

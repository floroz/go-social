package api

import "net/http"

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to write response")
	}
}

package api

import "net/http"

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

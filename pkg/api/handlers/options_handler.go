package handlers

import "net/http"

// HandleOptions is a generic handler for preflight OPTIONS requests
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebhookRoutes(r *mux.Router) {
	r.HandleFunc("/webhooks/{store_id}", WebhookHandler).Methods("POST")
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

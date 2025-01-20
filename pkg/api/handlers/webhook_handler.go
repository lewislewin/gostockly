package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterWebhookRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/webhooks/{store_id}", func(w http.ResponseWriter, r *http.Request) {
		WebhookHandler(w, r, db)
	}).Methods("POST")
}

func WebhookHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Use db to process webhook
	w.WriteHeader(http.StatusOK)
}

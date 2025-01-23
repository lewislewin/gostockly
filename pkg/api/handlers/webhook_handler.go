package handlers

import (
	"gostockly/internal/services"
	"gostockly/pkg/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type WebhookHandler struct {
	WebhookService *services.WebhookService
}

func RegisterWebhookRoutes(r *mux.Router, service *services.WebhookService) {
	handler := &WebhookHandler{WebhookService: service}
	r.HandleFunc("/webhook/orders", handler.HandleOrderWebhook).Methods("POST")
	r.HandleFunc("/webhook/products", handler.HandleProductWebhook).Methods("POST")
}

func (h *WebhookHandler) HandleOrderWebhook(w http.ResponseWriter, r *http.Request) {
	storeID := r.Header.Get("X-Store-ID") // Assuming the store ID is passed in the headers
	if storeID == "" {
		http.Error(w, "Missing store ID", http.StatusBadRequest)
		return
	}

	payload, err := utils.ReadRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = h.WebhookService.ProcessOrderWebhook(storeID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) HandleProductWebhook(w http.ResponseWriter, r *http.Request) {
	storeID := r.Header.Get("X-Store-ID") // Assuming the store ID is passed in the headers
	if storeID == "" {
		http.Error(w, "Missing store ID", http.StatusBadRequest)
		return
	}

	payload, err := utils.ReadRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = h.WebhookService.ProcessProductWebhook(storeID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

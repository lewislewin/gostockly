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
	r.HandleFunc("/webhook/orders", HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/webhook/orders", handler.HandleOrderWebhook).Methods("POST")
	r.HandleFunc("/webhook/products", handler.HandleProductWebhook).Methods("POST")
}

func (h *WebhookHandler) HandleOrderWebhook(w http.ResponseWriter, r *http.Request) {
	shopDomain := r.Header.Get("X-Shopify-Shop-Domain")
	if shopDomain == "" {
		http.Error(w, "Missing X-Shopify-Shop-Domain header", http.StatusBadRequest)
		return
	}

	payload, err := utils.ReadRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = h.WebhookService.ProcessOrderWebhook(shopDomain, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) HandleProductWebhook(w http.ResponseWriter, r *http.Request) {
	shopDomain := r.Header.Get("X-Shopify-Shop-Domain")
	if shopDomain == "" {
		http.Error(w, "Missing X-Shopify-Shop-Domain header", http.StatusBadRequest)
		return
	}

	payload, err := utils.ReadRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = h.WebhookService.ProcessProductWebhook(shopDomain, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

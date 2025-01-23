package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type InventoryHandler struct {
	InventoryService *services.InventoryService
}

func RegisterInventoryRoutes(r *mux.Router, service *services.InventoryService) {
	handler := &InventoryHandler{InventoryService: service}
	r.HandleFunc("/inventory/decrement", handler.DecrementInventory).Methods("POST")
}

func (h *InventoryHandler) DecrementInventory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKUs       []string `json:"skus"`
		Amount     int      `json:"amount"`
		StoreID    string   `json:"store_id"`
		LocationID string   `json:"location_id"`
		APIKey     string   `json:"api_key"`
		APISecret  string   `json:"api_secret"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	if len(req.SKUs) == 1 {
		err := h.InventoryService.DecrementSingleSKU(req.SKUs[0], req.Amount, storeID, req.LocationID, req.APIKey, req.APISecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := h.InventoryService.DecrementBulkSKUs(req.SKUs, req.Amount, storeID, req.LocationID, req.APIKey, req.APISecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

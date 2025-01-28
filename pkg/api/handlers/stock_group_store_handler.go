package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type StockGroupStoreHandler struct {
	StockGroupStoreService *services.StockGroupStoreService
}

func NewStockGroupStoreHandler(service *services.StockGroupStoreService) *StockGroupStoreHandler {
	return &StockGroupStoreHandler{StockGroupStoreService: service}
}

func (h *StockGroupStoreHandler) AddStoreToStockGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StockGroupID string `json:"stock_group_id"`
		StoreID      string `json:"store_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	stockGroupID, err := uuid.Parse(req.StockGroupID)
	if err != nil {
		http.Error(w, "Invalid stock group ID", http.StatusBadRequest)
		return
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	err = h.StockGroupStoreService.AddStoreToStockGroup(stockGroupID, storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Store added to stock group successfully"}`))
}

func RegisterStockGroupStoreRoutes(r *mux.Router, service *services.StockGroupStoreService) {
	handler := NewStockGroupStoreHandler(service)
	r.HandleFunc("/stockgroupstore/add", HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/stockgroupstore/add", handler.AddStoreToStockGroup).Methods("POST")
}

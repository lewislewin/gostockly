package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type StockGroupHandler struct {
	StockGroupService *services.StockGroupService
}

func NewStockGroupHandler(service *services.StockGroupService) *StockGroupHandler {
	return &StockGroupHandler{StockGroupService: service}
}

// CreateStockGroup handles creating a new stock group.
func (h *StockGroupHandler) CreateStockGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		CompanyID string `json:"company_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.CompanyID == "" {
		http.Error(w, "Name and CompanyID are required", http.StatusBadRequest)
		return
	}

	stockGroup, err := h.StockGroupService.CreateStockGroup(req.Name, req.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stockGroup)
}

// GetStockGroupsByCompany handles retrieving all stock groups for a company.
func (h *StockGroupHandler) GetStockGroupsByCompany(w http.ResponseWriter, r *http.Request) {
	companyID := mux.Vars(r)["company_id"]

	if _, err := uuid.Parse(companyID); err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	stockGroups, err := h.StockGroupService.GetStockGroupsByCompany(companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stockGroups)
}

// RegisterStockGroupRoutes registers stock group routes.
func RegisterStockGroupRoutes(r *mux.Router, service *services.StockGroupService) {
	handler := NewStockGroupHandler(service)
	r.HandleFunc("/stockgroup", handler.CreateStockGroup).Methods("POST")
	r.HandleFunc("/stockgroup/company/{company_id}", handler.GetStockGroupsByCompany).Methods("GET")
}

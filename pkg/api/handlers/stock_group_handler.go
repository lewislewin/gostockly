package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"gostockly/pkg/logger"
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
	log := logger.GetLogger()
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Error("Error: missing company_id in context")
		return
	}

	stockGroup, err := h.StockGroupService.CreateStockGroup(req.Name, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stockGroup)
}

// GetStockGroupsByCompany handles retrieving all stock groups for a company.
func (h *StockGroupHandler) GetStockGroupsByCompany(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Error("Error: missing company_id in context")
		return
	}

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
	r.HandleFunc("/stockgroups", HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/stockgroups", handler.CreateStockGroup).Methods("POST")
	r.HandleFunc("/stockgroups", HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/stockgroups", handler.GetStockGroupsByCompany).Methods("GET")
}

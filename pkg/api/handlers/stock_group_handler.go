package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"gostockly/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterStockGroupRoutes(r *mux.Router, stockGroupService *services.StockGroupService) {
	stockGroupRouter := r.PathPrefix("/stockgroups").Subrouter()

	stockGroupRouter.HandleFunc("", HandleOptions).Methods(http.MethodOptions)
	stockGroupRouter.HandleFunc("/{id}", HandleOptions).Methods(http.MethodOptions)

	stockGroupRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		CreateStockGroup(w, r, stockGroupService)
	}).Methods(http.MethodPost)

	stockGroupRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		GetStockGroupsByCompany(w, r, stockGroupService)
	}).Methods(http.MethodGet)

	stockGroupRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetStockGroupByID(w, r, stockGroupService)
	}).Methods(http.MethodGet)

	stockGroupRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateStockGroup(w, r, stockGroupService)
	}).Methods(http.MethodPut)

	stockGroupRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteStockGroup(w, r, stockGroupService)
	}).Methods(http.MethodDelete)
}

func CreateStockGroup(w http.ResponseWriter, r *http.Request, stockGroupService *services.StockGroupService) {
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

	stockGroup, err := stockGroupService.CreateStockGroup(req.Name, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stockGroup)
}

func GetStockGroupsByCompany(w http.ResponseWriter, r *http.Request, stockGroupService *services.StockGroupService) {
	log := logger.GetLogger()

	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Error("Error: missing company_id in context")
		return
	}

	stockGroups, err := stockGroupService.GetStockGroupsByCompany(companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stockGroups)
}

func GetStockGroupByID(w http.ResponseWriter, r *http.Request, stockGroupService *services.StockGroupService) {
	stockGroupID := mux.Vars(r)["id"]

	stockGroup, err := stockGroupService.GetStockGroupByID(stockGroupID)
	if err != nil {
		http.Error(w, "Stock group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stockGroup)
}

func UpdateStockGroup(w http.ResponseWriter, r *http.Request, stockGroupService *services.StockGroupService) {
	stockGroupID := mux.Vars(r)["id"]

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedStockGroup, err := stockGroupService.UpdateStockGroup(stockGroupID, req.Name)
	if err != nil {
		http.Error(w, "Failed to update stock group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStockGroup)
}

func DeleteStockGroup(w http.ResponseWriter, r *http.Request, stockGroupService *services.StockGroupService) {
	stockGroupID := mux.Vars(r)["id"]

	if err := stockGroupService.DeleteStockGroup(stockGroupID); err != nil {
		http.Error(w, "Failed to delete stock group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

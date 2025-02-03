package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterStoreRoutes(r *mux.Router, storeService *services.StoreService) {
	storeRouter := r.PathPrefix("/stores").Subrouter()

	// Allow OPTIONS requests for both list and single store routes
	storeRouter.HandleFunc("", HandleOptions).Methods(http.MethodOptions)
	storeRouter.HandleFunc("/{id}", HandleOptions).Methods(http.MethodOptions)

	// Store operations
	storeRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		CreateStore(w, r, storeService)
	}).Methods(http.MethodPost)

	storeRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ListStores(w, r, storeService)
	}).Methods(http.MethodGet)

	storeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetStoreByID(w, r, storeService)
	}).Methods(http.MethodGet)

	storeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateStore(w, r, storeService)
	}).Methods(http.MethodPut)

	storeRouter.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteStore(w, r, storeService)
	}).Methods(http.MethodDelete)
}

func CreateStore(w http.ResponseWriter, r *http.Request, storeService *services.StoreService) {
	var req struct {
		ShopifyStoreStub string `json:"shopify_store_stub"`
		AccessToken      string `json:"access_token"`
		WebhookSignature string `json:"webhook_signature"`
		LocationID       string `json:"location_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Println("Error: missing company_id in context")
		return
	}

	store, err := storeService.CreateStore(companyID, req.ShopifyStoreStub, req.AccessToken, req.WebhookSignature, req.LocationID)
	if err != nil {
		http.Error(w, "Failed to create store", http.StatusInternalServerError)
		log.Printf("Error creating store: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(store)
}

func ListStores(w http.ResponseWriter, r *http.Request, storeService *services.StoreService) {
	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Println("Error: missing company_id in context")
		return
	}

	stores, err := storeService.GetStoresByCompany(companyID)
	if err != nil {
		http.Error(w, "Failed to retrieve stores", http.StatusInternalServerError)
		log.Printf("Error retrieving stores: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stores)
}

func GetStoreByID(w http.ResponseWriter, r *http.Request, storeService *services.StoreService) {
	storeID := mux.Vars(r)["id"]

	store, err := storeService.GetStoreByID(storeID)
	if err != nil {
		http.Error(w, "Store not found", http.StatusNotFound)
		log.Printf("Error retrieving store: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(store)
}

func UpdateStore(w http.ResponseWriter, r *http.Request, storeService *services.StoreService) {
	storeID := mux.Vars(r)["id"]

	var req struct {
		ShopifyStoreStub string `json:"shopify_store_stub"`
		AccessToken      string `json:"access_token"`
		WebhookSignature string `json:"webhook_signature"`
		LocationID       string `json:"location_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	updatedStore, err := storeService.UpdateStore(storeID, req.ShopifyStoreStub, req.AccessToken, req.WebhookSignature, req.LocationID)
	if err != nil {
		http.Error(w, "Failed to update store", http.StatusInternalServerError)
		log.Printf("Error updating store: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updatedStore)
}

func DeleteStore(w http.ResponseWriter, r *http.Request, storeService *services.StoreService) {
	storeID := mux.Vars(r)["id"]

	if err := storeService.DeleteStore(storeID); err != nil {
		http.Error(w, "Failed to delete store", http.StatusInternalServerError)
		log.Printf("Error deleting store: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

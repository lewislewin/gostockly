package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterStoreRoutes(r *mux.Router, storeService *services.StoreService) {
	r.HandleFunc("/stores", HandleOptions).Methods(http.MethodOptions)
	r.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		CreateStore(w, r, storeService)
	}).Methods("POST")

	r.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		ListStores(w, r, storeService)
	}).Methods("GET")
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

	// Get company ID from context
	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Println("Error: missing company_id in context")
		return
	}

	// Call the service to create the store
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
	// Get company ID from context
	companyID, ok := r.Context().Value("company_id").(string)
	if !ok || companyID == "" {
		http.Error(w, "Unauthorized: missing company_id in context", http.StatusUnauthorized)
		log.Println("Error: missing company_id in context")
		return
	}

	// Fetch stores for the company
	stores, err := storeService.GetStoresByCompany(companyID)
	if err != nil {
		http.Error(w, "Failed to retrieve stores", http.StatusInternalServerError)
		log.Printf("Error retrieving stores: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stores)
}

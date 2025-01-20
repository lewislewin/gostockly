package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterInventoryRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/inventory/update", func(w http.ResponseWriter, r *http.Request) {
		UpdateInventory(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/inventory/{sku}", func(w http.ResponseWriter, r *http.Request) {
		GetInventory(w, r, db)
	}).Methods("GET")
}

func UpdateInventory(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Use db to update inventory
	w.WriteHeader(http.StatusOK)
}

func GetInventory(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Use db to retrieve inventory
	w.WriteHeader(http.StatusOK)
}

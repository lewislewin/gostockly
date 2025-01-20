package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterInventoryRoutes(r *mux.Router) {
	r.HandleFunc("/inventory/update", UpdateInventory).Methods("POST")
	r.HandleFunc("/inventory/{sku}", GetInventory).Methods("GET")
}

func UpdateInventory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetInventory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

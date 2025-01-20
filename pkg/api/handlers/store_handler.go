package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterStoreRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		CreateStore(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		ListStores(w, r, db)
	}).Methods("GET")
}

func CreateStore(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Use the db to handle store creation
	w.WriteHeader(http.StatusCreated)
}

func ListStores(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Use the db to list stores
	w.WriteHeader(http.StatusOK)
}

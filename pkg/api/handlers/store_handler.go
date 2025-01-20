package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterStoreRoutes(r *mux.Router) {
	r.HandleFunc("/stores", CreateStore).Methods("POST")
	r.HandleFunc("/stores", ListStores).Methods("GET")
}

func CreateStore(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func ListStores(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

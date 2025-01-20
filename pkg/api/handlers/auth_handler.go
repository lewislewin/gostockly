package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", RegisterUser).Methods("POST")
	r.HandleFunc("/auth/login", LoginUser).Methods("POST")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"gostockly/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router, userService *services.UserService) {
	r.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterUser(w, r, userService)
	}).Methods("POST")

	r.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		LoginUser(w, r, userService)
	}).Methods("POST")
}

type RegisterUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
	Subdomain   string `json:"subdomain"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request, userService *services.UserService) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	user, err := userService.RegisterUser(req.Email, req.Password, req.CompanyName, req.Subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error registering user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(w http.ResponseWriter, r *http.Request, userService *services.UserService) {
	var req LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	token, err := userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("Error authenticating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

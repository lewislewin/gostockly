package handlers

import (
	"encoding/json"
	"gostockly/config"
	"gostockly/internal/repositories"
	"gostockly/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *mux.Router, cfg *config.Config, db *gorm.DB) {
	// Initialize the repositories and services
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	userService := services.NewUserService(userRepo, companyRepo, cfg.JWTSecret)

	// Register routes
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

	// Call the service
	user, err := userService.RegisterUser(req.Email, req.Password, req.CompanyName, req.Subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error registering user: %v", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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

	// Call the service
	token, err := userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("Error authenticating user: %v", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

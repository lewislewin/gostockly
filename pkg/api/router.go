package api

import (
	"gostockly/config"
	"gostockly/internal/repositories"
	"gostockly/internal/services"
	"gostockly/pkg/api/handlers"
	"gostockly/pkg/logger"
	"gostockly/pkg/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *mux.Router {
	// Initialize logger
	log := logger.GetLogger()

	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	storeRepo := repositories.NewStoreRepository(db)

	// Create services
	userService := services.NewUserService(userRepo, companyRepo, cfg.JWTSecret)
	storeService := services.NewStoreService(storeRepo)

	// Create router
	r := mux.NewRouter()

	// Apply global logging middleware
	r.Use(middleware.LoggingMiddleware)

	// Log route registration
	log.Info("Registering routes...")

	// Public routes
	handlers.RegisterAuthRoutes(r, userService)
	log.Info("Auth routes registered")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(userService)) // Auth middleware uses userService
	handlers.RegisterStoreRoutes(protected, storeService)
	log.Info("Store routes registered")

	log.Info("All routes registered successfully")
	return r
}

package api

import (
	"gostockly/config"
	"gostockly/pkg/api/handlers"
	"gostockly/pkg/logger"
	"gostockly/pkg/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *mux.Router {
	// Initialize logger
	log := logger.GetLogger()

	// Create a new router
	r := mux.NewRouter()

	// Apply global logging middleware
	r.Use(middleware.LoggingMiddleware)

	// Log route registration
	log.Info("Registering routes...")

	// Public routes (no authentication required)
	handlers.RegisterAuthRoutes(r, cfg, db)
	log.Info("Auth routes registered")

	// Protected routes (authentication required)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Register protected routes
	handlers.RegisterStoreRoutes(protected, db)
	log.Info("Store routes registered")

	handlers.RegisterInventoryRoutes(protected, db)
	log.Info("Inventory routes registered")

	handlers.RegisterWebhookRoutes(protected, db)
	log.Info("Webhook routes registered")

	log.Info("All routes registered successfully")
	return r
}

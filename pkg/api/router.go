package api

import (
	"gostockly/config"
	"gostockly/pkg/api/handlers"
	"gostockly/pkg/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	// Add middleware
	r.Use(middleware.LoggingMiddleware)

	// Register routes
	handlers.RegisterAuthRoutes(r)
	handlers.RegisterStoreRoutes(r)
	handlers.RegisterInventoryRoutes(r)
	handlers.RegisterWebhookRoutes(r)

	return r
}

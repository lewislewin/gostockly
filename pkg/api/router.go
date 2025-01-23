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
	log := logger.GetLogger()

	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	storeRepo := repositories.NewStoreRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	stockGroupStoreRepo := repositories.NewStockGroupStoreRepository(db)

	userService := services.NewUserService(userRepo, companyRepo, cfg.JWTSecret)
	storeService := services.NewStoreService(storeRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, storeRepo)
	webhookService := services.NewWebhookService(storeRepo, inventoryRepo, stockGroupStoreRepo)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	log.Info("Registering routes...")
	handlers.RegisterAuthRoutes(r, userService)
	log.Info("Auth routes registered")

	handlers.RegisterWebhookRoutes(r, webhookService)
	log.Info("Webhook routes registered")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(userService))
	handlers.RegisterStoreRoutes(protected, storeService)
	handlers.RegisterInventoryRoutes(protected, inventoryService)
	log.Info("Inventory routes registered")

	log.Info("All routes registered successfully")
	return r
}

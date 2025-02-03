package repositories

import (
	"errors"
	"gostockly/internal/models"
	"strings"

	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// CreateStore adds a new store to the database.
func (r *StoreRepository) CreateStore(store *models.Store) error {
	return r.db.Create(store).Error
}

// GetStoresByCompany retrieves all stores belonging to a specific company.
func (r *StoreRepository) GetStoresByCompany(companyID string) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Where("company_id = ?", companyID).Find(&stores).Error
	return stores, err
}

// GetStoreByID retrieves a store by its ID.
func (r *StoreRepository) GetStoreByID(storeID string) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, "id = ?", storeID).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) GetStoreByShopifyStub(shopifyStub string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("shopify_store_stub = ?", shopifyStub).First(&store).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &store, err
}

// GetStoreByShopifyDomain retrieves a store by its Shopify domain.
func (r *StoreRepository) GetStoreByShopifyDomain(domain string) (*models.Store, error) {
	// Extract the subdomain from the Shopify domain
	parts := strings.Split(domain, ".")
	if len(parts) < 3 || parts[1] != "myshopify" || parts[2] != "com" {
		return nil, errors.New("invalid Shopify domain")
	}
	shopifyStub := parts[0] // Subdomain is the first part of the domain

	// Use the subdomain to retrieve the store
	return r.GetStoreByShopifyStub(shopifyStub)
}

func (r *StoreRepository) UpdateStore(store *models.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) DeleteStore(storeID string) error {
	return r.db.Delete(&models.Store{}, "id = ?", storeID).Error
}

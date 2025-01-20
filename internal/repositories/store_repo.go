package repositories

import (
	"gostockly/internal/models"

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

// GetStoresByCompany retrieves all stores belonging to a company.
func (r *StoreRepository) GetStoresByCompany(companyID string) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Where("company_id = ?", companyID).Find(&stores).Error
	return stores, err
}

// GetStoresByID retrieves all stores belonging to a company.
func (r *StoreRepository) GetStoreByID(storeID string) (models.Store, error) {
	var store models.Store
	err := r.db.Where("store_id = ?", storeID).Find(&store).Error
	return store, err
}

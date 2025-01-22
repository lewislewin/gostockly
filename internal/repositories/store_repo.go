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

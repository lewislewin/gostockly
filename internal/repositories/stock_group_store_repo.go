package repositories

import (
	"errors"
	"gostockly/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockGroupStoreRepository struct {
	db *gorm.DB
}

func NewStockGroupStoreRepository(db *gorm.DB) *StockGroupStoreRepository {
	return &StockGroupStoreRepository{db: db}
}

// GetStoresByStockGroup retrieves all stores belonging to a specific stock group.
func (r *StockGroupStoreRepository) GetStoresByStockGroup(stockGroupID uuid.UUID) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Joins("JOIN stock_group_stores ON stock_group_stores.store_id = stores.id").
		Where("stock_group_stores.stock_group_id = ?", stockGroupID).
		Find(&stores).Error
	return stores, err
}

// GetStockGroupsByStore retrieves the stock group a store belongs to.
func (r *StockGroupStoreRepository) GetStockGroupsByStore(storeID uuid.UUID) (*models.StockGroup, error) {
	var stockGroup models.StockGroup
	err := r.db.Joins("JOIN stock_group_stores ON stock_group_stores.stock_group_id = stock_groups.id").
		Where("stock_group_stores.store_id = ?", storeID).
		First(&stockGroup).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &stockGroup, err
}

// AddStoreToStockGroup adds a store to a stock group.
func (r *StockGroupStoreRepository) AddStoreToStockGroup(stockGroupID, storeID uuid.UUID) error {
	// Ensure the store is not already part of another stock group
	existingStockGroup, err := r.GetStockGroupsByStore(storeID)
	if err != nil {
		return err
	}
	if existingStockGroup != nil {
		return errors.New("store is already part of a stock group")
	}

	stockGroupStore := &models.StockGroupStore{
		ID:           uuid.New(),
		StockGroupID: stockGroupID,
		StoreID:      storeID,
	}
	return r.db.Create(stockGroupStore).Error
}

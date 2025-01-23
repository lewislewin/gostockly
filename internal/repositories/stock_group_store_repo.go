package repositories

import (
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

// GetStockGroupsByStore retrieves all stock groups a store belongs to.
func (r *StockGroupStoreRepository) GetStockGroupsByStore(storeID uuid.UUID) ([]models.StockGroup, error) {
	var stockGroups []models.StockGroup
	err := r.db.Joins("JOIN stock_group_stores ON stock_group_stores.stock_group_id = stock_groups.id").
		Where("stock_group_stores.store_id = ?", storeID).
		Find(&stockGroups).Error
	return stockGroups, err
}

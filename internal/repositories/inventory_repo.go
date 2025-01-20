package repositories

import (
	"errors"
	"gostockly/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) CreateInventory(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetInventoryBySKU(sku string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("sku = ?", sku).First(&inventory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory not found")
	}
	return &inventory, err
}

func (r *InventoryRepository) UpdateStockLevel(sku string, stockLevel int) error {
	result := r.db.Model(&models.Inventory{}).Where("sku = ?", sku).Update("stock_level", stockLevel)
	if result.RowsAffected == 0 {
		return errors.New("no inventory found to update")
	}
	return result.Error
}

func (r *InventoryRepository) GetInventoryByStockGroup(stockGroup string) ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.Where("stock_group = ?", stockGroup).Find(&inventories).Error
	return inventories, err
}

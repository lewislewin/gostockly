package repositories

import (
	"errors"
	"gostockly/internal/models"

	"github.com/google/uuid"
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

func (r *InventoryRepository) GetInventoryBySKUAndStore(sku string, storeID uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("sku = ? AND store_id = ?", sku, storeID).First(&inventory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory not found for the specified store")
	}
	return &inventory, err
}

func (r *InventoryRepository) UpdateInventoryItemID(sku string, storeID uuid.UUID, inventoryItemID string) error {
	result := r.db.Model(&models.Inventory{}).Where("sku = ? AND store_id = ?", sku, storeID).
		Update("inventory_item_id", inventoryItemID)
	if result.RowsAffected == 0 {
		return errors.New("no inventory found to update for the specified store")
	}
	return result.Error
}

func (r *InventoryRepository) GetInventoryByStockGroupAndStore(stockGroup string, storeID uuid.UUID) ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.Where("stock_group = ? AND store_id = ?", stockGroup, storeID).Find(&inventories).Error
	return inventories, err
}

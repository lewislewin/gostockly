package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	SKU             string    `gorm:"not null;" json:"sku"`
	StockGroup      string    `gorm:"not null;" json:"stock_group"`
	InventoryItemID string    `gorm:"not null;" json:"inventory_item_id"`
	StoreID         uuid.UUID `gorm:"not null;" json:"store_id"`
	UpdatedAt       time.Time `json:"updated_at"`
}

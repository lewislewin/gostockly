package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	SKU        string    `gorm:"not null;unique" json:"sku"`
	StockLevel int       `gorm:"not null" json:"stock_level"`
	StockGroup string    `gorm:"not null" json:"stock_group"`
	UpdatedAt  time.Time `json:"updated_at"`
}

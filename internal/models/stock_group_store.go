package models

import (
	"time"

	"github.com/google/uuid"
)

type StockGroupStore struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StockGroupID uuid.UUID `gorm:"type:uuid;not null" json:"stock_group_id"`
	StoreID      uuid.UUID `gorm:"type:uuid;not null" json:"store_id"`
	CreatedAt    time.Time `json:"created_at"`

	StockGroup *StockGroup `gorm:"foreignKey:StockGroupID" json:"stock_group"`
	Store      *Store      `gorm:"foreignKey:StoreID" json:"store"`
}

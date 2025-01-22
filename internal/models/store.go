package models

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyID        uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	ShopifyStoreStub string    `gorm:"not null" json:"shopify_store_stub"`
	APIKey           string    `gorm:"not null" json:"api_key"`
	APISecret        string    `gorm:"not null" json:"api_secret"`
	CreatedAt        time.Time `json:"created_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company"`
}

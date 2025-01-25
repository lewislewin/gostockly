package models

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyID        uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	ShopifyStoreStub string    `gorm:"not null" json:"shopify_store_stub"`
	AccessToken      string    `gorm:"not null" json:"access_token"`
	WebhookSignature string    `gorm:"not null;default:''" json:"webhook_signature"` // Default empty string
	LocationID       string    `gorm:"not null;default:''" json:"location_id"`       // Default empty string
	CreatedAt        time.Time `json:"created_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company"`
}

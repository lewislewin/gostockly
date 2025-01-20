package models

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Subdomain string    `gorm:"unique;not null" json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`

	Users       []User       `gorm:"foreignKey:CompanyID" json:"users"`
	Stores      []Store      `gorm:"foreignKey:CompanyID" json:"stores"`
	StockGroups []StockGroup `gorm:"foreignKey:CompanyID" json:"stock_groups"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt time.Time `json:"created_at"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company"`
}

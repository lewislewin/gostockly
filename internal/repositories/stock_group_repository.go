package repositories

import (
	"gostockly/internal/models"

	"gorm.io/gorm"
)

type StockGroupRepository struct {
	db *gorm.DB
}

func NewStockGroupRepository(db *gorm.DB) *StockGroupRepository {
	return &StockGroupRepository{db: db}
}

// CreateStockGroup adds a new stock group to the database.
func (r *StockGroupRepository) CreateStockGroup(stockGroup *models.StockGroup) error {
	return r.db.Create(stockGroup).Error
}

// GetStockGroupsByCompany retrieves all stock groups belonging to a company.
func (r *StockGroupRepository) GetStockGroupsByCompany(companyID string) ([]models.StockGroup, error) {
	var stockGroups []models.StockGroup
	err := r.db.Where("company_id = ?", companyID).Find(&stockGroups).Error
	return stockGroups, err
}

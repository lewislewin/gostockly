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

func (r *StockGroupRepository) CreateStockGroup(stockGroup *models.StockGroup) error {
	return r.db.Create(stockGroup).Error
}

func (r *StockGroupRepository) GetStockGroupsByCompany(companyID string) ([]models.StockGroup, error) {
	var stockGroups []models.StockGroup
	err := r.db.Where("company_id = ?", companyID).Find(&stockGroups).Error
	return stockGroups, err
}

func (r *StockGroupRepository) GetStockGroupByID(stockGroupID string) (*models.StockGroup, error) {
	var stockGroup models.StockGroup
	err := r.db.First(&stockGroup, "id = ?", stockGroupID).Error
	if err != nil {
		return nil, err
	}
	return &stockGroup, nil
}

func (r *StockGroupRepository) UpdateStockGroup(stockGroup *models.StockGroup) error {
	return r.db.Save(stockGroup).Error
}

func (r *StockGroupRepository) DeleteStockGroup(stockGroupID string) error {
	return r.db.Delete(&models.StockGroup{}, "id = ?", stockGroupID).Error
}

package services

import (
	"errors"
	"gostockly/internal/models"
	"gostockly/internal/repositories"

	"github.com/google/uuid"
)

type StockGroupService struct {
	StockGroupRepo *repositories.StockGroupRepository
}

func NewStockGroupService(stockGroupRepo *repositories.StockGroupRepository) *StockGroupService {
	return &StockGroupService{
		StockGroupRepo: stockGroupRepo,
	}
}

// CreateStockGroup creates a new stock group for a company.
func (s *StockGroupService) CreateStockGroup(name string, companyID string) (*models.StockGroup, error) {
	companyUUID, err := uuid.Parse(companyID)
	if err != nil {
		return nil, errors.New("invalid company ID")
	}

	stockGroup := &models.StockGroup{
		ID:        uuid.New(),
		Name:      name,
		CompanyID: companyUUID,
	}

	err = s.StockGroupRepo.CreateStockGroup(stockGroup)
	if err != nil {
		return nil, errors.New("failed to create stock group")
	}

	return stockGroup, nil
}

// GetStockGroupsByCompany retrieves all stock groups belonging to a company.
func (s *StockGroupService) GetStockGroupsByCompany(companyID string) ([]models.StockGroup, error) {
	return s.StockGroupRepo.GetStockGroupsByCompany(companyID)
}

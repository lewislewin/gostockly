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
	return &StockGroupService{StockGroupRepo: stockGroupRepo}
}

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

func (s *StockGroupService) GetStockGroupsByCompany(companyID string) ([]models.StockGroup, error) {
	return s.StockGroupRepo.GetStockGroupsByCompany(companyID)
}

func (s *StockGroupService) GetStockGroupByID(stockGroupID string) (*models.StockGroup, error) {
	return s.StockGroupRepo.GetStockGroupByID(stockGroupID)
}

func (s *StockGroupService) UpdateStockGroup(stockGroupID, name string) (*models.StockGroup, error) {
	stockGroup, err := s.StockGroupRepo.GetStockGroupByID(stockGroupID)
	if err != nil {
		return nil, err
	}

	stockGroup.Name = name

	if err := s.StockGroupRepo.UpdateStockGroup(stockGroup); err != nil {
		return nil, err
	}

	return stockGroup, nil
}

func (s *StockGroupService) DeleteStockGroup(stockGroupID string) error {
	return s.StockGroupRepo.DeleteStockGroup(stockGroupID)
}

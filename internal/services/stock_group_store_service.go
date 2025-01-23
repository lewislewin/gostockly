package services

import (
	"errors"
	"gostockly/internal/repositories"

	"github.com/google/uuid"
)

type StockGroupStoreService struct {
	StockGroupStoreRepo *repositories.StockGroupStoreRepository
	StockGroupRepo      *repositories.StockGroupRepository
	StoreRepo           *repositories.StoreRepository
}

func NewStockGroupStoreService(
	stockGroupStoreRepo *repositories.StockGroupStoreRepository,
	stockGroupRepo *repositories.StockGroupRepository,
	storeRepo *repositories.StoreRepository,
) *StockGroupStoreService {
	return &StockGroupStoreService{
		StockGroupStoreRepo: stockGroupStoreRepo,
		StockGroupRepo:      stockGroupRepo,
		StoreRepo:           storeRepo,
	}
}

// AddStoreToStockGroup adds a store to a stock group.
func (s *StockGroupStoreService) AddStoreToStockGroup(stockGroupID, storeID uuid.UUID) error {
	// Ensure the stock group exists
	stockGroup, err := s.StockGroupRepo.GetStockGroupsByCompany(stockGroupID.String())
	if err != nil || stockGroup == nil {
		return errors.New("stock group not found")
	}

	// Ensure the store exists
	_, err = s.StoreRepo.GetStoreByID(storeID.String())
	if err != nil {
		return errors.New("store not found")
	}

	// Add the store to the stock group
	return s.StockGroupStoreRepo.AddStoreToStockGroup(stockGroupID, storeID)
}

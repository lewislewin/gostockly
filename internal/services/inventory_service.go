package services

import (
	"errors"
	"gostockly/internal/repositories"
)

type InventoryService struct {
	InventoryRepo *repositories.InventoryRepository
	StoreRepo     *repositories.StoreRepository
}

func NewInventoryService(inventoryRepo *repositories.InventoryRepository, storeRepo *repositories.StoreRepository) *InventoryService {
	return &InventoryService{
		InventoryRepo: inventoryRepo,
		StoreRepo:     storeRepo,
	}
}

// UpdateStock updates the stock level for a specific SKU in a stock group.
func (s *InventoryService) UpdateStock(sku string, stockLevel int, stockGroup string) error {
	inventory, err := s.InventoryRepo.GetInventoryBySKU(sku)
	if err != nil {
		return err
	}

	if inventory.StockGroup != stockGroup {
		return errors.New("SKU does not belong to the specified stock group")
	}

	return s.InventoryRepo.UpdateStockLevel(sku, stockLevel)
}

// SyncStockGroup ensures all stores in a stock group have the same stock level for a SKU.
func (s *InventoryService) SyncStockGroup(sku string, stockLevel int, stockGroup string) error {
	inventories, err := s.InventoryRepo.GetInventoryByStockGroup(stockGroup)
	if err != nil {
		return err
	}

	for _, inventory := range inventories {
		if inventory.SKU == sku {
			err := s.InventoryRepo.UpdateStockLevel(inventory.SKU, stockLevel)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

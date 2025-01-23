package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"gostockly/internal/repositories"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
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

// ShopifyRequest sends a GraphQL request to Shopify
func (s *InventoryService) ShopifyRequest(apiKey, apiSecret, query string, variables map[string]interface{}) ([]byte, error) {
	url := "https://your-shopify-store.myshopify.com/admin/api/2023-01/graphql.json"

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	return ioutil.ReadAll(resp.Body)
}

// DecrementSingleSKU decrements inventory for a specific SKU and store
func (s *InventoryService) DecrementSingleSKU(sku string, amount int, storeID uuid.UUID, locationID, apiKey, apiSecret string) error {
	inventory, err := s.InventoryRepo.GetInventoryBySKUAndStore(sku, storeID)
	if err != nil {
		return err
	}

	query := `
		mutation AdjustInventoryLevel($inventoryItemId: ID!, $locationId: ID!, $adjustment: Int!) {
			inventoryAdjustQuantity(
				inventoryItemId: $inventoryItemId,
				locationId: $locationId,
				availableDelta: $adjustment
			) {
				inventoryLevel {
					id
					available
				}
				userErrors {
					field
					message
				}
			}
		}
	`

	variables := map[string]interface{}{
		"inventoryItemId": inventory.InventoryItemID,
		"locationId":      locationID,
		"adjustment":      -amount,
	}

	_, err = s.ShopifyRequest(apiKey, apiSecret, query, variables)
	return err
}

// DecrementBulkSKUs decrements inventory for multiple SKUs for a specific store
func (s *InventoryService) DecrementBulkSKUs(skus []string, amount int, storeID uuid.UUID, locationID, apiKey, apiSecret string) error {
	for _, sku := range skus {
		err := s.DecrementSingleSKU(sku, amount, storeID, locationID, apiKey, apiSecret)
		if err != nil {
			return err
		}
	}
	return nil
}

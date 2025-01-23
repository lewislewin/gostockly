package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/shopify"

	"github.com/google/uuid"
)

type WebhookService struct {
	StoreRepo           *repositories.StoreRepository
	InventoryRepo       *repositories.InventoryRepository
	StockGroupStoreRepo *repositories.StockGroupStoreRepository
}

func NewWebhookService(
	storeRepo *repositories.StoreRepository,
	inventoryRepo *repositories.InventoryRepository,
	stockGroupStoreRepo *repositories.StockGroupStoreRepository,
) *WebhookService {
	return &WebhookService{
		StoreRepo:           storeRepo,
		InventoryRepo:       inventoryRepo,
		StockGroupStoreRepo: stockGroupStoreRepo,
	}
}

// ProcessOrderWebhook processes an order webhook from Shopify and updates stock across the stock group.
func (s *WebhookService) ProcessOrderWebhook(storeID string, payload []byte) error {
	// Validate store existence
	sourceStore, err := s.StoreRepo.GetStoreByID(storeID)
	if err != nil {
		return errors.New("invalid store ID")
	}

	// Parse the webhook payload
	var order struct {
		LineItems []struct {
			SKU      string `json:"sku"`
			Quantity int    `json:"quantity"`
		} `json:"line_items"`
	}
	err = json.Unmarshal(payload, &order)
	if err != nil {
		return errors.New("failed to parse webhook payload")
	}

	// Get the stock group for the source store
	stockGroups, err := s.StockGroupStoreRepo.GetStockGroupsByStore(sourceStore.ID)
	if err != nil || len(stockGroups) == 0 {
		return errors.New("no stock group found for this store")
	}
	stockGroup := stockGroups[0] // A store belongs to only one stock group

	// Get all stores in the stock group
	stores, err := s.StockGroupStoreRepo.GetStoresByStockGroup(stockGroup.ID)
	if err != nil {
		return errors.New("failed to retrieve stores in the stock group")
	}

	// Iterate through stores in the stock group
	for _, targetStore := range stores {
		// Skip the source store
		if targetStore.ID == sourceStore.ID {
			continue
		}

		shopifyClient := shopify.NewShopifyClient(targetStore.APIKey, targetStore.APISecret, targetStore.ShopifyStoreStub)

		// Collect bulk adjustments
		var adjustments []map[string]interface{}
		for _, item := range order.LineItems {
			inventory, err := s.InventoryRepo.GetInventoryBySKUAndStore(item.SKU, targetStore.ID)
			if err != nil {
				continue // Log or handle SKU not found
			}

			adjustments = append(adjustments, map[string]interface{}{
				"inventoryItemId": inventory.InventoryItemID,
				"locationId":      "location-id-for-target-store", // Replace with actual location ID
				"adjustment":      -item.Quantity,
			})

			// Send in batches of 250
			if len(adjustments) == 250 {
				err := s.sendBulkInventoryAdjustment(shopifyClient, adjustments)
				if err != nil {
					// Log error, continue processing remaining adjustments
					continue
				}
				adjustments = nil // Reset for the next batch
			}
		}

		// Send remaining adjustments
		if len(adjustments) > 0 {
			err := s.sendBulkInventoryAdjustment(shopifyClient, adjustments)
			if err != nil {
				// Log error
				continue
			}
		}
	}

	return nil
}

// ProcessProductWebhook processes a product creation/update webhook from Shopify.
func (s *WebhookService) ProcessProductWebhook(storeID string, payload []byte) error {
	// Validate store existence
	store, err := s.StoreRepo.GetStoreByID(storeID)
	if err != nil {
		return errors.New("invalid store ID")
	}

	// Parse the webhook payload
	var product struct {
		Variants []struct {
			SKU             string `json:"sku"`
			InventoryItemID string `json:"inventory_item_id"`
		} `json:"variants"`
	}
	err = json.Unmarshal(payload, &product)
	if err != nil {
		return errors.New("failed to parse product webhook payload")
	}

	// Update the database with new or updated SKUs and InventoryItemIDs
	for _, variant := range product.Variants {
		if variant.SKU == "" || variant.InventoryItemID == "" {
			continue // Skip invalid variants
		}

		// Update or create the inventory record
		_, err := s.InventoryRepo.GetInventoryBySKUAndStore(variant.SKU, store.ID)
		if err != nil {
			// If not found, create a new record
			newInventory := &models.Inventory{
				ID:              uuid.New(),
				SKU:             variant.SKU,
				InventoryItemID: variant.InventoryItemID,
				StoreID:         store.ID,
			}
			err := s.InventoryRepo.CreateInventory(newInventory)
			if err != nil {
				return errors.New("failed to create inventory: " + err.Error())
			}
		} else {
			// Update existing record
			err := s.InventoryRepo.UpdateInventoryItemID(variant.SKU, store.ID, variant.InventoryItemID)
			if err != nil {
				return errors.New("failed to update inventory: " + err.Error())
			}
		}
	}

	return nil
}

func (s *WebhookService) sendBulkInventoryAdjustment(client *shopify.ShopifyClient, adjustments []map[string]interface{}) error {
	// Build the GraphQL mutation
	query := "mutation AdjustInventoryLevels {"
	for i, adjustment := range adjustments {
		query += `
			item` + string(rune(i)) + `: inventoryAdjustQuantity(
				inventoryItemId: "` + adjustment["inventoryItemId"].(string) + `",
				locationId: "` + adjustment["locationId"].(string) + `",
				availableDelta: ` + fmt.Sprintf("%d", adjustment["adjustment"].(int)) + `
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
		`
	}
	query += "}"

	// Send the request
	_, err := client.SendGraphQLRequest(query, nil)
	return err
}

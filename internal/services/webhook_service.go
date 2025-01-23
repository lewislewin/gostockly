package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/logger"
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
func (s *WebhookService) ProcessOrderWebhook(shopDomain string, payload []byte) error {
	log := logger.GetLogger()
	log.Info("Processing order webhook for shop: %s", shopDomain)

	// Validate store existence
	sourceStore, err := s.StoreRepo.GetStoreByShopifyDomain(shopDomain)
	if err != nil {
		log.Error("Failed to find store for shop domain %s: %v", shopDomain, err)
		return errors.New("invalid store domain")
	}
	log.Info("Found source store: %s (ID: %s)", sourceStore.ShopifyStoreStub, sourceStore.ID)

	// Parse the webhook payload
	var order struct {
		LineItems []struct {
			SKU      string `json:"sku"`
			Quantity int    `json:"quantity"`
		} `json:"line_items"`
	}
	err = json.Unmarshal(payload, &order)
	if err != nil {
		log.Error("Failed to parse webhook payload for shop %s: %v", shopDomain, err)
		return errors.New("failed to parse webhook payload")
	}
	log.Info("Parsed %d line items for shop %s", len(order.LineItems), shopDomain)

	// Get the stock group for the source store
	stockGroup, err := s.StockGroupStoreRepo.GetStockGroupsByStore(sourceStore.ID)
	if err != nil {
		log.Error("No stock group found for store %s: %v", sourceStore.ID, err)
		return errors.New("no stock group found for this store")
	}
	log.Info("Found stock group: %s for store %s", stockGroup.ID, sourceStore.ID)

	// Get all stores in the stock group
	stores, err := s.StockGroupStoreRepo.GetStoresByStockGroup(stockGroup.ID)
	if err != nil {
		log.Error("Failed to retrieve stores in stock group %s: %v", stockGroup.ID, err)
		return errors.New("failed to retrieve stores in the stock group")
	}
	log.Info("Found %d stores in stock group %s", len(stores), stockGroup.ID)

	// Iterate through stores in the stock group
	for _, targetStore := range stores {
		if targetStore.ID == sourceStore.ID {
			continue
		}
		log.Info("Processing stock updates for target store: %s (ID: %s)", targetStore.ShopifyStoreStub, targetStore.ID)

		shopifyClient := shopify.NewShopifyClient(targetStore.AccessToken, targetStore.ShopifyStoreStub)

		// Process each line item individually
		for _, item := range order.LineItems {
			inventory, err := s.InventoryRepo.GetInventoryBySKUAndStore(item.SKU, targetStore.ID)
			if err != nil {
				log.Debug("Failed to find inventory for SKU %s in store %s: %v", item.SKU, targetStore.ID, err)
				continue
			}

			adjustment := map[string]interface{}{
				"inventoryItemId": inventory.InventoryItemID,
				"locationId":      "105135735110", // Replace with actual location ID
				"adjustment":      -item.Quantity,
			}

			// Send the adjustment
			log.Info("Sending inventory adjustment for SKU: %s to store: %s", item.SKU, targetStore.ShopifyStoreStub)
			err = s.sendInventoryAdjustment(shopifyClient, adjustment)
			if err != nil {
				log.Error("Failed to send inventory adjustment for store %s: %v", targetStore.ShopifyStoreStub, err)
			}
		}
	}

	log.Info("Finished processing order webhook for shop: %s", shopDomain)
	return nil
}

// ProcessProductWebhook processes a product creation/update webhook from Shopify.
func (s *WebhookService) ProcessProductWebhook(shopDomain string, payload []byte) error {
	log := logger.GetLogger()
	log.Info("Processing product webhook for shop: %s", shopDomain)

	// Validate store existence
	store, err := s.StoreRepo.GetStoreByShopifyDomain(shopDomain)
	if err != nil {
		log.Error("Failed to find store for shop domain %s: %v", shopDomain, err)
		return errors.New("invalid store domain")
	}
	log.Info("Found store: %s (ID: %s)", store.ShopifyStoreStub, store.ID)

	// Parse the webhook payload
	var product struct {
		Variants []struct {
			SKU             string `json:"sku"`
			InventoryItemID string `json:"inventory_item_id"`
		} `json:"variants"`
	}
	err = json.Unmarshal(payload, &product)
	if err != nil {
		log.Error("Failed to parse product webhook payload for shop %s: %v", shopDomain, err)
		return errors.New("failed to parse product webhook payload")
	}
	log.Info("Parsed %d variants for shop %s", len(product.Variants), shopDomain)

	// Update the database with new or updated SKUs and InventoryItemIDs
	for _, variant := range product.Variants {
		if variant.SKU == "" || variant.InventoryItemID == "" {
			log.Debug("Skipping variant with missing SKU or InventoryItemID for shop %s", shopDomain)
			continue
		}

		_, err := s.InventoryRepo.GetInventoryBySKUAndStore(variant.SKU, store.ID)
		if err != nil {
			log.Info("Creating new inventory for SKU %s in store %s", variant.SKU, store.ID)
			newInventory := &models.Inventory{
				ID:              uuid.New(),
				SKU:             variant.SKU,
				InventoryItemID: variant.InventoryItemID,
				StoreID:         store.ID,
			}
			err := s.InventoryRepo.CreateInventory(newInventory)
			if err != nil {
				log.Error("Failed to create inventory for SKU %s in store %s: %v", variant.SKU, store.ID, err)
			}
		} else {
			log.Info("Updating inventory for SKU %s in store %s", variant.SKU, store.ID)
			err := s.InventoryRepo.UpdateInventoryItemID(variant.SKU, store.ID, variant.InventoryItemID)
			if err != nil {
				log.Error("Failed to update inventory for SKU %s in store %s: %v", variant.SKU, store.ID, err)
			}
		}
	}

	log.Info("Finished processing product webhook for shop: %s", shopDomain)
	return nil
}

func (s *WebhookService) sendInventoryAdjustment(client *shopify.ShopifyClient, adjustment map[string]interface{}) error {
	log := logger.GetLogger()

	// Define the GraphQL mutation
	query := `
		mutation inventoryAdjustQuantities($input: InventoryAdjustQuantitiesInput!) {
			inventoryAdjustQuantities(input: $input) {
				userErrors {
					field
					message
				}
				inventoryAdjustmentGroup {
					createdAt
					reason
					referenceDocumentUri
					changes {
						name
						delta
					}
				}
			}
		}
	`

	// Ensure inventoryItemId and locationId are formatted as global IDs
	inventoryItemId := fmt.Sprintf("gid://shopify/InventoryItem/%s", adjustment["inventoryItemId"])
	locationId := fmt.Sprintf("gid://shopify/Location/%s", adjustment["locationId"])

	// Construct the mutation variables
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"reason":               "correction",
			"name":                 "available",
			"referenceDocumentUri": "logistics://some.warehouse/take/2023-01/13",
			"changes": []map[string]interface{}{
				{
					"delta":           adjustment["adjustment"],
					"inventoryItemId": inventoryItemId,
					"locationId":      locationId,
				},
			},
		},
	}

	// Log the query and variables for debugging
	log.Debug("Query: %s", query)
	log.Debug("Variables: %+v", variables)

	// Send the GraphQL request
	body, err := client.SendGraphQLRequest(query, variables)
	if err != nil {
		log.Error("Failed to send inventory adjustment mutation: %v", err)
		return err
	}

	// Log the response for debugging
	log.Debug("Response body: %s", body)

	// Parse the response to check for errors
	var response struct {
		Data struct {
			InventoryAdjustQuantities struct {
				UserErrors []struct {
					Field   string `json:"field"`
					Message string `json:"message"`
				} `json:"userErrors"`
			} `json:"inventoryAdjustQuantities"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Error("Failed to parse response from Shopify: %v", err)
		return err
	}

	// Log user errors if any are present
	if len(response.Data.InventoryAdjustQuantities.UserErrors) > 0 {
		for _, userError := range response.Data.InventoryAdjustQuantities.UserErrors {
			log.Error("Shopify user error: field=%s, message=%s", userError.Field, userError.Message)
		}
		return errors.New("Shopify returned user errors for inventory adjustment")
	}

	log.Info("Successfully sent inventory adjustment mutation to Shopify")
	return nil
}

package services

import (
	"encoding/json"
	"errors"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/utils"
	"net/http"
)

type WebhookService struct {
	StoreRepo     *repositories.StoreRepository
	InventoryRepo *repositories.InventoryRepository
}

func NewWebhookService(storeRepo *repositories.StoreRepository, inventoryRepo *repositories.InventoryRepository) *WebhookService {
	return &WebhookService{
		StoreRepo:     storeRepo,
		InventoryRepo: inventoryRepo,
	}
}

// ProcessWebhook processes an incoming Shopify webhook payload.
func (s *WebhookService) ProcessWebhook(storeID string, payload []byte) error {
	// Validate store existence
	_, err := s.StoreRepo.GetStoreByID(storeID)
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

	// Deduct stock for each SKU in the order
	for _, item := range order.LineItems {
		inventory, err := s.InventoryRepo.GetInventoryBySKU(item.SKU)
		if err != nil {
			return err
		}

		newStock := inventory.StockLevel - item.Quantity
		if newStock < 0 {
			return errors.New("insufficient stock for SKU: " + item.SKU)
		}

		err = s.InventoryRepo.UpdateStockLevel(item.SKU, newStock)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateWebhook ensures the webhook payload is authentic.
func (s *WebhookService) ValidateWebhook(req *http.Request, store *models.Store) error {
	hmacHeader := req.Header.Get("X-Shopify-Hmac-SHA256")
	if hmacHeader == "" {
		return errors.New("missing HMAC header")
	}

	// Validate HMAC using the store's secret
	payload, err := utils.ReadRequestBody(req)
	if err != nil {
		return errors.New("failed to read request body")
	}

	valid := utils.ValidateHMAC(payload, hmacHeader, store.APISecret)
	if !valid {
		return errors.New("invalid HMAC signature")
	}

	return nil
}

package main

import (
	"fmt"
	"gostockly/internal/models"
	"gostockly/internal/repositories"
	"gostockly/pkg/database"
	"gostockly/pkg/shopify"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	// Connect to the database
	db := database.Connect()

	// Initialize repositories
	storeRepo := repositories.NewStoreRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)

	// Get company ID from arguments or environment variable
	companyID := os.Getenv("COMPANY_ID")
	if companyID == "" {
		log.Fatal("COMPANY_ID is not set. Pass it as an environment variable.")
	}

	// Fetch all stores for the company
	stores, err := storeRepo.GetStoresByCompany(companyID)
	if err != nil {
		log.Fatalf("Failed to fetch stores: %v", err)
	}

	// Process each store
	for _, store := range stores {
		log.Printf("Fetching products for store: %s", store.ShopifyStoreStub)

		shopifyClient := shopify.NewShopifyClient(store.AccessToken, store.ShopifyStoreStub)
		products, err := shopifyClient.FetchProducts()
		if err != nil {
			log.Printf("Failed to fetch products for store %s: %v", store.ShopifyStoreStub, err)
			continue
		}

		// Sync inventory for the store
		for _, product := range products {
			for _, variant := range product.Variants {
				err := syncInventory(inventoryRepo, store, variant)
				if err != nil {
					log.Printf("Failed to sync inventory for SKU %s in store %s: %v", variant.SKU, store.ShopifyStoreStub, err)
				}
			}
		}
	}

	log.Println("Finished syncing products for all stores")
}

func syncInventory(repo *repositories.InventoryRepository, store models.Store, variant shopify.Variant) error {
	// Check if inventory already exists
	_, err := repo.GetInventoryBySKUAndStore(variant.SKU, store.ID)
	if err != nil {
		// Create new inventory if it doesn't exist
		log.Printf("Creating new inventory for SKU %s in store %s", variant.SKU, store.ShopifyStoreStub)
		newInventory := &models.Inventory{
			ID:              uuid.New(),
			SKU:             variant.SKU,
			InventoryItemID: fmt.Sprintf("%d", variant.InventoryItemID),
			StoreID:         store.ID,
		}
		return repo.CreateInventory(newInventory)
	}

	// Update existing inventory if necessary
	log.Printf("Updating inventory for SKU %s in store %s", variant.SKU, store.ShopifyStoreStub)
	return repo.UpdateInventoryItemID(variant.SKU, store.ID, fmt.Sprintf("%d", variant.InventoryItemID))
}

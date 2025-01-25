package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gostockly/pkg/logger"
	"io"
	"net/http"
)

type ShopifyClient struct {
	AccessToken string
	StoreURL    string // Full Shopify API URL
}

func NewShopifyClient(accessToken, storeStub string) *ShopifyClient {
	// Ensure the store URL is formatted correctly
	storeURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2025-01", storeStub)
	log := logger.GetLogger()
	log.Info("Creating Shopify client for store: %s", storeStub)

	return &ShopifyClient{
		AccessToken: accessToken,
		StoreURL:    storeURL,
	}
}

func (c *ShopifyClient) SendGraphQLRequest(query string, variables map[string]interface{}) ([]byte, error) {
	log := logger.GetLogger()

	// Log request details
	log.Debug("Sending GraphQL request to Shopify: %s", c.StoreURL+"/graphql.json")
	log.Debug("Query: %s", query)
	if variables != nil {
		variablesJSON, _ := json.Marshal(variables)
		log.Debug("Variables: %s", variablesJSON)
	}

	// Build the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		log.Error("Failed to marshal request body: %v", err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", c.StoreURL+"/graphql.json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Failed to create HTTP request: %v", err)
		return nil, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", c.AccessToken)
	log.Debug("Headers added to request: Content-Type=application/json, X-Shopify-Access-Token=[REDACTED]")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send request to Shopify: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Log response status
	log.Info("Received response from Shopify: status=%d", resp.StatusCode)

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("Shopify API request failed with status: %s, body: %s", resp.Status, body)
		return nil, fmt.Errorf("Shopify API request failed with status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body: %v", err)
		return nil, err
	}

	log.Debug("Response body: %s", body)
	return body, nil
}

// Product represents a Shopify product.
type Product struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Variants []Variant `json:"variants"`
}

// Variant represents a Shopify product variant.
type Variant struct {
	ID                int64  `json:"id"`
	SKU               string `json:"sku"`
	InventoryItemID   int64  `json:"inventory_item_id"`
	InventoryQuantity int    `json:"inventory_quantity"`
}

// FetchProducts retrieves all products from the Shopify store.
func (c *ShopifyClient) FetchProducts() ([]Product, error) {
	apiURL := c.StoreURL + "/products.json"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Shopify-Access-Token", c.AccessToken)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Products []Product `json:"products"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Products, nil
}

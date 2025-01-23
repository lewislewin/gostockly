package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"gostockly/pkg/logger"
	"io"
	"net/http"
)

type ShopifyClient struct {
	AccessToken string
	StoreURL    string
}

func NewShopifyClient(accessToken, storeStub string) *ShopifyClient {
	// Ensure the store URL is formatted correctly
	storeURL := "https://" + storeStub + ".myshopify.com/admin/api/2025-01/graphql.json"
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
	log.Debug("Sending GraphQL request to Shopify: %s", c.StoreURL)
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
	req, err := http.NewRequest("POST", c.StoreURL, bytes.NewBuffer(requestBody))
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
		log.Error("Shopify API request failed with status: %s", resp.Status)
		return nil, errors.New("Shopify API request failed with status: " + resp.Status)
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

package shopify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ShopifyClient struct {
	APIKey    string
	APISecret string
	StoreURL  string
}

func NewShopifyClient(apiKey, apiSecret, storeStub string) *ShopifyClient {
	return &ShopifyClient{
		APIKey:    apiKey,
		APISecret: apiSecret,
		StoreURL:  "https://" + storeStub + ".myshopify.com/admin/api/2023-01/graphql.json",
	}
}

func (c *ShopifyClient) SendGraphQLRequest(query string, variables map[string]interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.StoreURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.APIKey, c.APISecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Shopify API request failed with status: " + resp.Status)
	}

	return io.ReadAll(resp.Body)
}

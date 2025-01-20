package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ValidateHMAC validates the HMAC signature of a payload against a given secret.
func ValidateHMAC(payload []byte, hmacHeader, secret string) bool {
	// Create a new HMAC using SHA256 and the provided secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)

	// Encode the calculated HMAC in Base64 to match Shopify's header
	expectedHMAC := base64.StdEncoding.EncodeToString(expectedMAC)

	// Compare the calculated HMAC with the provided HMAC header
	return hmacHeader == expectedHMAC
}

package utils

import (
	"bytes"
	"io"
	"net/http"
)

// ReadRequestBody reads the body of an HTTP request and returns it as a byte slice.
func ReadRequestBody(req *http.Request) ([]byte, error) {
	// Read the body into a byte slice
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// Reset the body so it can be read again downstream
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

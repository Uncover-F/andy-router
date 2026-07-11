package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"
)

func VerifyAndyKey(key string) (bool, error) {
	// Define API request
	reqBody := []byte(`{
		"model": "auto",
		"messages": [
			{
				"role": "user",
				"content": "ping"
			}
		],
		"max_tokens": 1
	}`)
	req, err := http.NewRequest("POST", "https://andy.mindcraft-ce.com/api/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	// Send API request
	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.New("unable to read response body: " + err.Error())
	}

	// Validate API request
	if resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden {
		return false, errors.New("invalid Andy API key: " + string(body))
	}

	return true, nil
}

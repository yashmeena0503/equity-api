package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *APIClient) SendMessage(req MessageRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshall request: %w", err)
	}

	resp, err := c.Client.Post(
		c.BaseURL+"/message",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("API returned %s", resp.Status)
	}

	return nil
}

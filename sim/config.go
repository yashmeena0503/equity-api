package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseURL   string              `json:"base_url"`
	Platforms map[string][]string `json:"platforms"`
}

func NewConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

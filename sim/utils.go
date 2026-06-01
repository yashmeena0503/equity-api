package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	Workers int `json:"workers"`
	RequestPerWorker int `json:"request_per_worker"`
	DelayBetweenRequestsMS int `json:"delay_between_requests_ms"`
	PlatformNames []string `json:"platform_names"`
	MaxMessageLength int `json:"max_message_length"`

	DelayBetweenRequests time.Duration `json:"-"`
}
func newConfig(filename string) (Config, error) {
	var config Config
	// Read external JSON file.
	bytes, err := os.ReadFile(filename); if err != nil { return config, err }
	// Unmarshal (convert) JSON bytes into the struct
	err = json.Unmarshal(bytes, &config); if err != nil { return config, err }
	// Convert miliseconds integer into proper time.Duration
	config.DelayBetweenRequests = time.Duration(config.DelayBetweenRequestsMS) * time.Millisecond
	return config, nil
}

type endpoint struct {
	method string
	path string
	body func(c *Config) any // nil = no body
}

// Helpers.

func randomPlatformName(c *Config) string {
	return c.PlatformNames[rand.Intn(len(c.PlatformNames))]
}

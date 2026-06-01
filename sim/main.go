package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// Endpoints.
var endpoints = []endpoint {
	// {
	// 	method: "GET",
	// 	path: "/payout",
	// },
	{
		method: "POST",
		path: "/message",
		body: func(c *Config) any {
			return map[string] any {
				"platform_name": randomPlatformName(c),
				"platform_id": rand.Intn(c.Workers),
				"message_length": rand.Intn(c.MaxMessageLength),
			}
		},
	},
}

func main() {

	// Initialize Env.
	port := os.Getenv("PORT"); if port == "" { port = "6767" }
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	slog.Info("Initialized Env")

	// Initialize Config.
	config, err := newConfig("config.json"); if err != nil { panic(err) }

	println("Starting Simulation.")

	var wg sync.WaitGroup
	results := make(chan result, config.Workers * config.RequestPerWorker)
	for i := range config.Workers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(baseURL, &config, workerID, results)
		}(i)
	}
	// Close results channel once all workers finish.
	go func() {
		wg.Wait()
		close(results)
	}()
	summarise(baseURL, results)

}

func worker(baseURL string, c *Config, id int, out chan <- result) {
	client := &http.Client { Timeout: 5 * time.Second }
	for range c.RequestPerWorker {
		ep := endpoints[rand.Intn(len(endpoints))]
		res := fire(baseURL, c, client, ep)
		res.worker = id
		out <- res
		// random sleep so we don't thunderherd.
		time.Sleep(time.Duration(rand.Int63n(int64(c.DelayBetweenRequests))))
	}
}

// HTTP fire

type result struct {
	worker int
	method string
	path string
	statusCode int
	err error
	duration time.Duration
}

func fire(baseURL string, c *Config, client *http.Client, ep endpoint) result {
	r := result { method: ep.method, path: ep.path }
	var buf *bytes.Buffer
	if ep.body != nil {
		data, _ := json.Marshal(ep.body(c))
		buf = bytes.NewBuffer(data)
	} else {
		buf = &bytes.Buffer {}
	}
	req, err := http.NewRequest(ep.method, baseURL + ep.path, buf)
	if err != nil { r.err = err; return r }
	req.Header.Set("Content-Type", "application/json")
	start := time.Now()
	resp, err := client.Do(req)
	r.duration = time.Since(start)
	if err != nil { r.err = err; return r }
	defer resp.Body.Close()
	r.statusCode = resp.StatusCode
	return r
}

func summarise(baseURL string, results <- chan result) {
	counts := map[int]int {}
	errors := 0
	total := 0
	var totalDur time.Duration

	for r := range results {
		total ++
		totalDur += r.duration
		if r.err != nil {
			errors++
			log.Printf("[worker %d] ERROR %s %s -> %v", r.worker, r.method, r.path, r.err)
		} else {
			counts[r.statusCode]++
			log.Printf("[worker %d] %s %s -> %d (%s)", r.worker, r.method, r.path, r.statusCode, r.duration.Round(time.Millisecond))
		}
	}
	fmt.Println("\n=== Summary ===")
	fmt.Printf("URL: %s\n", baseURL)
	fmt.Printf("Total Requests: %d\n", total)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Avg Latency: %s\n", (totalDur / time.Duration(max(total, 1))).Round(time.Millisecond))
	fmt.Println("Status codes:")
	for code, count := range counts { fmt.Printf("  %d -> %d\n", code, count) }

}
package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// SyncClient is an adapter that implements ports.SyncClient
type SyncClient struct {
	timeout time.Duration
	verbose bool
	client  *http.Client
}

// NewSyncClient creates a new SyncClient
func NewSyncClient(timeout time.Duration, verbose bool) *SyncClient {
	return &SyncClient{
		timeout: timeout,
		verbose: verbose,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// TriggerSync sends a POST request to the sync endpoint
func (sc *SyncClient) TriggerSync(ctx context.Context, url string, apiKey string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("X-Admin-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoJob/1.0")

	if sc.verbose {
		log.Printf("📤 Sending POST request to: %s", url)
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if sc.verbose {
		log.Printf("📥 Response status: %d", resp.StatusCode)
		if len(body) > 0 {
			log.Printf("   Body: %s", string(body))
		}
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	return resp.StatusCode, nil
}

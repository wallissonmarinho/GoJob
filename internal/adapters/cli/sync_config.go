package cli

import "time"

// SyncConfig holds configuration for the sync operation
type SyncConfig struct {
	URL     string
	APIKey  string
	Timeout time.Duration
	Verbose bool
}

// NewSyncConfig creates a new SyncConfig
func NewSyncConfig(url string, apiKey string, timeout time.Duration, verbose bool) *SyncConfig {
	return &SyncConfig{
		URL:     url,
		APIKey:  apiKey,
		Timeout: timeout,
		Verbose: verbose,
	}
}

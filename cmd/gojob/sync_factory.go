package main

import (
	"fmt"
	"time"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoJob/internal/adapters/cli"
	"github.com/wallissonmarinho/GoJob/internal/adapters/commands"
	"github.com/wallissonmarinho/GoJob/internal/adapters/http"
	"github.com/wallissonmarinho/GoJob/internal/app/config"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
	"github.com/wallissonmarinho/GoJob/internal/core/services"
)

// BuildSyncExecutorFactory returns a commands.ExecutorFactory which constructs
// the Sync executor (wires HTTP client, sync service, and config).
// It uses the loaded config and allows CLI flags to override values.
func BuildSyncExecutorFactory(cfg config.Config) commands.ExecutorFactory {
	return func(c *urfavecli.Context) (ports.Executor, error) {
		// Get values from flags, fallback to config
		url := c.String("url")
		if url == "" {
			url = cfg.SyncURL
		}

		apiKey := c.String("api-key")
		if apiKey == "" {
			apiKey = cfg.APIKey
		}

		// CLI flags override config timeout
		timeout := cfg.Timeout
		if c.IsSet("timeout") {
			timeout = time.Duration(c.Int("timeout")) * time.Second
		}

		// CLI flags override config verbose
		verbose := cfg.Verbose
		if c.IsSet("verbose") {
			verbose = c.Bool("verbose")
		}

		if url == "" {
			return nil, fmt.Errorf("sync URL is required (--url or SYNC_URL environment variable)")
		}

		if apiKey == "" {
			return nil, fmt.Errorf("API key is required (--api-key or API_KEY environment variable)")
		}

		// Create CLI config
		cliConfig := cli.NewSyncConfig(url, apiKey, timeout, verbose)

		// Adapter: HTTP client
		httpAdapter := http.NewSyncClient(timeout, verbose)

		// Service: Sync orchestration
		syncService := services.NewSyncService(httpAdapter, cliConfig)

		return syncService, nil
	}
}

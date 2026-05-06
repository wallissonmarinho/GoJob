package main

import (
	"fmt"
	"time"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoJob/internal/adapters/cli"
	"github.com/wallissonmarinho/GoJob/internal/adapters/commands"
	"github.com/wallissonmarinho/GoJob/internal/adapters/http"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
	"github.com/wallissonmarinho/GoJob/internal/core/services"
)

// BuildSyncExecutorFactory returns a commands.ExecutorFactory which constructs
// the Sync executor (wires HTTP client, sync service, and config).
func BuildSyncExecutorFactory() commands.ExecutorFactory {
	return func(c *urfavecli.Context) (ports.Executor, error) {
		url := c.String("url")
		apiKey := c.String("api-key")
		timeout := time.Duration(c.Int("timeout")) * time.Second
		verbose := c.Bool("verbose")

		if url == "" {
			return nil, fmt.Errorf("sync URL is required (--url or SYNC_URL)")
		}

		if apiKey == "" {
			return nil, fmt.Errorf("API key is required (--api-key or API_KEY)")
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

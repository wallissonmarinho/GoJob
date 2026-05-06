package main

import (
	"log"
	"os"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoJob/internal/app/config"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Build commands using factory pattern with a provided executor factory
	commandFactory := BuildCommandFactory(cfg)

	app := &urfavecli.App{
		Name:        "gojob",
		Usage:       "HTTP Sync Trigger - Universal cronjob for triggering sync endpoints",
		Description: "Generic CLI tool to trigger HTTP sync endpoints, designed to run as a Kubernetes CronJob",
		Commands:    commandFactory.BuildCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
}

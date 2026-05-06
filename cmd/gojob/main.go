package main

import (
	"log"
	"os"

	urfavecli "github.com/urfave/cli/v2"
)

func main() {
	// Build commands using factory pattern with a provided executor factory
	commandFactory := BuildCommandFactory()

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

package commands

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v2"
)

// SyncCommandHandler handles the sync CLI command
type SyncCommandHandler struct {
	executorFactory ExecutorFactory
}

// NewSyncCommandHandler creates a new SyncCommandHandler
func NewSyncCommandHandler(executorFactory ExecutorFactory) *SyncCommandHandler {
	return &SyncCommandHandler{
		executorFactory: executorFactory,
	}
}

// BuildCommand returns the sync CLI command
func (sch *SyncCommandHandler) BuildCommand() *urfavecli.Command {
	return &urfavecli.Command{
		Name:    "sync",
		Aliases: []string{"s"},
		Usage:   "Trigger HTTP sync endpoint",
		Flags: []urfavecli.Flag{
			&urfavecli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "HTTP sync endpoint URL",
				EnvVars: []string{"SYNC_URL"},
			},
			&urfavecli.StringFlag{
				Name:    "api-key",
				Aliases: []string{"k"},
				Usage:   "API key for authentication",
				EnvVars: []string{"API_KEY"},
			},
			&urfavecli.IntFlag{
				Name:  "timeout",
				Value: 30,
				Usage: "Request timeout in seconds",
			},
			&urfavecli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging",
			},
		},
		Action: func(c *urfavecli.Context) error {
			executor, err := sch.executorFactory(c)
			if err != nil {
				return fmt.Errorf("failed to create executor: %w", err)
			}

			command := NewSyncCommand(executor)
			return command.Execute()
		},
	}
}

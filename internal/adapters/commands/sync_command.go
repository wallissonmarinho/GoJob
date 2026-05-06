package commands

import "github.com/wallissonmarinho/GoJob/internal/core/ports"

// SyncCommand implements ports.ServiceCommand
type SyncCommand struct {
	executor ports.Executor
}

// NewSyncCommand creates a new SyncCommand
func NewSyncCommand(executor ports.Executor) ports.ServiceCommand {
	return &SyncCommand{executor: executor}
}

// Execute runs the executor
func (sc *SyncCommand) Execute() error {
	return sc.executor.Execute()
}

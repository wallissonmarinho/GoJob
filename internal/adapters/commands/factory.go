package commands

import (
	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
)

// ExecutorFactory is a factory function that creates an Executor from CLI context
type ExecutorFactory func(c *urfavecli.Context) (ports.Executor, error)

// CommandFactory builds all available CLI commands
type CommandFactory struct {
	syncHandler *SyncCommandHandler
}

// NewCommandFactoryWithExecutorFactories creates a CommandFactory with an executor factory
func NewCommandFactoryWithExecutorFactories(syncFactory ExecutorFactory) *CommandFactory {
	return &CommandFactory{
		syncHandler: NewSyncCommandHandler(syncFactory),
	}
}

// NewCommandFactory is a convenience constructor used in tests or simple setups
func NewCommandFactory() *CommandFactory {
	dummy := func(c *urfavecli.Context) (ports.Executor, error) { return nil, nil }
	return NewCommandFactoryWithExecutorFactories(dummy)
}

// BuildCommands returns all available CLI commands
func (f *CommandFactory) BuildCommands() []*urfavecli.Command {
	cmds := []*urfavecli.Command{}

	if f.syncHandler != nil {
		cmds = append(cmds, f.syncHandler.BuildCommand())
	}

	return cmds
}

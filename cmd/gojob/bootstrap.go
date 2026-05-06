package main

import (
	"github.com/wallissonmarinho/GoJob/internal/adapters/commands"
	"github.com/wallissonmarinho/GoJob/internal/app/config"
)

// BuildCommandFactory composes and returns a CommandFactory wired with the
// executor factory. Keeps composition in one place for easier maintenance.
func BuildCommandFactory(cfg config.Config) *commands.CommandFactory {
	return commands.NewCommandFactoryWithExecutorFactories(
		BuildSyncExecutorFactory(cfg),
	)
}

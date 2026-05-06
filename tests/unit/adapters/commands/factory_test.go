package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	urfavecli "github.com/urfave/cli/v2"
	adaptcmd "github.com/wallissonmarinho/GoJob/internal/adapters/commands"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
)

func TestCommandFactory_BuildCommands(t *testing.T) {
	dummyFactory := func(c *urfavecli.Context) (ports.Executor, error) {
		return nil, nil
	}

	factory := adaptcmd.NewCommandFactoryWithExecutorFactories(dummyFactory)
	commands := factory.BuildCommands()

	assert.NotNil(t, commands)
	assert.Len(t, commands, 1)
	assert.Equal(t, "sync", commands[0].Name)
}

func TestCommandFactory_BuildCommands_HasAliases(t *testing.T) {
	dummyFactory := func(c *urfavecli.Context) (ports.Executor, error) {
		return nil, nil
	}

	factory := adaptcmd.NewCommandFactoryWithExecutorFactories(dummyFactory)
	commands := factory.BuildCommands()

	assert.NotNil(t, commands[0].Aliases)
	assert.Contains(t, commands[0].Aliases, "s")
}

func TestSyncCommandHandler_BuildCommand_Flags(t *testing.T) {
	dummyFactory := func(c *urfavecli.Context) (ports.Executor, error) {
		return nil, nil
	}

	handler := adaptcmd.NewSyncCommandHandler(dummyFactory)
	command := handler.BuildCommand()

	assert.Equal(t, "sync", command.Name)
	assert.NotEmpty(t, command.Flags)

	// Check for expected flags
	flagNames := map[string]bool{}
	for _, flag := range command.Flags {
		switch v := flag.(type) {
		case *urfavecli.StringFlag:
			flagNames[v.Name] = true
		case *urfavecli.IntFlag:
			flagNames[v.Name] = true
		case *urfavecli.BoolFlag:
			flagNames[v.Name] = true
		}
	}

	assert.True(t, flagNames["url"])
	assert.True(t, flagNames["api-key"])
	assert.True(t, flagNames["timeout"])
	assert.True(t, flagNames["verbose"])
}

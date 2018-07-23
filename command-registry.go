package godatabus

import (
	"sync"
	"fmt"
)

type CommandHandler func(command Command) error

// CommandHandlerRegistry stores the handlers for commands
type CommandHandlerRegistry interface {
	Add(command Command, handler CommandHandler)
	Get(command Command) (CommandHandler, error)
}

// CommandRegistry contains a registry of command-handler style
type CommandRegistry struct {
	sync.RWMutex
	registry map[string]CommandHandler
	resolver NameResolver
}

// NewCommandRegistry creates a new CommandHandler
func NewCommandRegistry(resolver NameResolver) *CommandRegistry {
	return &CommandRegistry{
		registry: make(map[string]CommandHandler),
		resolver: resolver,
	}
}

// Add a new command with its handler
func (c *CommandRegistry) Add(command Command, handler CommandHandler) {
	c.Lock()
	defer c.Unlock()

	name := c.resolver.Add(command)
	c.registry[name] = handler
}

// Get the handler for a command
func (c *CommandRegistry) Get(command Command) (CommandHandler, error) {
	name, _ := c.resolver.ReverseResolve(Message(command.(Message)))

	handler, ok := c.registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}
	return handler, nil
}

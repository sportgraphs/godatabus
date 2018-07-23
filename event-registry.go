package godatabus

import (
	"sync"
	"fmt"
)

type EventHandler func(event Event)

// CommandHandlerRegistry stores the handlers for commands
type EventHandlerRegistry interface {
	Add(event Event, handler EventHandler)
	Get(event Event) ([]EventHandler, error)
}

// CommandRegistry contains a registry of command-handler style
type EventRegistry struct {
	sync.RWMutex
	registry map[string][]EventHandler
	resolver NameResolver
}

// NewCommandRegistry creates a new CommandHandler
func NewEventRegistry(resolver NameResolver) *EventRegistry {
	return &EventRegistry{
		registry: make(map[string][]EventHandler),
		resolver: resolver,
	}
}

// Add a new command with its handler
func (c *EventRegistry) Add(event Event, handler EventHandler) {
	c.Lock()
	defer c.Unlock()

	name := c.resolver.Add(event)
	c.registry[name] = append(c.registry[name], handler)
}

// Get the handler for a command
func (c *EventRegistry) Get(event Event) ([]EventHandler, error) {
	name, _ := c.resolver.ReverseResolve(event)

	handler, ok := c.registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}
	return handler, nil
}

package godatabus

import (
	"reflect"
	"sync"
	"fmt"
)

type NameResolver interface {
	Add(message Message) string
	Resolve(name string) (reflect.Type, error)
	ReverseResolve(message Message) (string, error)
}

type TypeNameResolver struct {
	sync.RWMutex
	registry map[string]reflect.Type
}

func NewTypeNameResolver() *TypeNameResolver {
	return &TypeNameResolver{
		registry: make(map[string]reflect.Type),
	}
}

func (resolver *TypeNameResolver) Add(message Message) string {
	resolver.Lock()
	defer resolver.Unlock()

	resolver.registry[message.GetType()] = reflect.TypeOf(message)

	return message.GetType()
}

func (resolver *TypeNameResolver) Resolve(name string) (reflect.Type, error) {
	messageType, ok := resolver.registry[name]

	if !ok {
		return nil, fmt.Errorf("can't resolve message type for %s", name)
	}

	return messageType, nil
}

func (resolver *TypeNameResolver) ReverseResolve(message Message) (string, error) {
	return message.GetType(), nil
}

package godatabus

import (
	"sync"
	"context"
)

type MessageBus interface {
	Handle(ctx context.Context, message Message) error
}

type MessageBusSupportingMiddleware struct {
	sync.RWMutex
	middlewares []Middleware
}

// NewCommandRegistry creates a new CommandHandler
func NewMessageBusSupportingMiddleware() *MessageBusSupportingMiddleware {
	return &MessageBusSupportingMiddleware{
		middlewares: make([]Middleware, 0),
	}
}

func (messageBus *MessageBusSupportingMiddleware) AppendMiddleware(middleware Middleware) {
	messageBus.Lock()
	defer messageBus.Unlock()

	messageBus.middlewares = append(messageBus.middlewares, middleware)
}


func (messageBus *MessageBusSupportingMiddleware) Handle(ctx context.Context, message Message) (context.Context, error) {
	return messageBus.callableForNextMiddleware(0)(ctx, message)
}


func (messageBus *MessageBusSupportingMiddleware) callableForNextMiddleware(index int) MiddlewareHandler {
	indexExists := len(messageBus.middlewares) > index

	if false == indexExists {
		return func(ctx context.Context, message Message) (context.Context, error) {return ctx, nil}
	}
	middleware := messageBus.middlewares[index]

	return func(ctx context.Context, message Message) (context.Context, error) {
		return middleware.Handle(ctx, message, messageBus.callableForNextMiddleware(index + 1))
	}
}

package godatabus

import "sync"

type MessageBus interface {
	Handle(message Message) error
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


func (messageBus *MessageBusSupportingMiddleware) Handle(message Message) error {
	return messageBus.callableForNextMiddleware(0)(message)
}


func (messageBus *MessageBusSupportingMiddleware) callableForNextMiddleware(index int) MiddlewareHandler {
	indexExists := len(messageBus.middlewares) > index

	if false == indexExists {
		return func(message Message) error {return nil}
	}
	middleware := messageBus.middlewares[index]

	return func(message Message) error {
		return middleware.Handle(message, messageBus.callableForNextMiddleware(index + 1))
	}
}

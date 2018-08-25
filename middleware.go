package godatabus

import (
	"github.com/sportgraphs/rabbit"
	"encoding/json"
	"github.com/streadway/amqp"
	"context"
)

type MiddlewareHandler func(ctx context.Context, message Message) (context.Context, error)

type Middleware interface {
	Handle(ctx context.Context, message Message, next MiddlewareHandler) (context.Context, error)
}

type DelegatesMessageToHandlerMiddleware struct {
	commandHandlerRegistry CommandHandlerRegistry
}

func NewDelegatesMessageToHandlerMiddleware(registry CommandHandlerRegistry) *DelegatesMessageToHandlerMiddleware {
	return &DelegatesMessageToHandlerMiddleware{
		commandHandlerRegistry: registry,
	}
}

func (middleware *DelegatesMessageToHandlerMiddleware) Handle(ctx context.Context, message Message, next MiddlewareHandler) (context.Context, error) {
	handler, err := middleware.commandHandlerRegistry.Get(Command(message))

	if err != nil {
		return ctx, err
	}

	err = handler(ctx, Command(message))
	if err != nil {
		return ctx, err
	}

	return next(ctx, message)
}

type AlwaysPublishesMessagesMiddleware struct {
	mq                 rabbit.MQ
	producerName       string
	nameResolver       NameResolver
	routingKeyResolver RoutingKeyResolver
}

func NewAlwaysPublishesMessagesMiddleware(producerName string, mq rabbit.MQ, nameResolver NameResolver, routingKeyResolver RoutingKeyResolver) *AlwaysPublishesMessagesMiddleware {
	return &AlwaysPublishesMessagesMiddleware{
		mq:                 mq,
		producerName:       producerName,
		nameResolver:       nameResolver,
		routingKeyResolver: routingKeyResolver,
	}
}

type NotifiesMessageSubscribersMiddleware struct {
	eventHandlerRegistry EventHandlerRegistry
}

func NewNotifiesMessageSubscribersMiddleware(registry EventHandlerRegistry) *NotifiesMessageSubscribersMiddleware {
	return &NotifiesMessageSubscribersMiddleware{
		eventHandlerRegistry: registry,
	}
}

func (middleware *NotifiesMessageSubscribersMiddleware) Handle(ctx context.Context, message Message, next MiddlewareHandler) (context.Context, error) {
	subscribers, err := middleware.eventHandlerRegistry.Get(Event(message))

	if err != nil {
		return next(ctx, message)
	}

	for _, subscriber := range subscribers {
		err = subscriber(ctx, Event(message))

		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, message)
}

func (middleware *AlwaysPublishesMessagesMiddleware) Handle(ctx context.Context, message Message, next MiddlewareHandler) (context.Context, error) {
	err := middleware.publish(message)

	if err != nil {
		return ctx, err
	}

	return next(ctx, message)
}

func (middleware *AlwaysPublishesMessagesMiddleware) publish(message Message) error {
	producer, err := middleware.mq.GetPublisher(middleware.producerName)

	if err != nil {
		return err
	}

	messageType, err := middleware.nameResolver.ReverseResolve(message)
	if err != nil {
		return err
	}

	routingKey := middleware.routingKeyResolver.Resolve(message)
	serializedMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	envelope := Envelope{
		Type:    messageType,
		Payload: string(serializedMessage),
	}

	payload, err := json.Marshal(&envelope)

	if err != nil {
		return err
	}

	producer.PublishWithRoutingKey(
		amqp.Publishing{
			Body: payload,
		}, routingKey)

	return nil
}

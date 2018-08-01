package godatabus

import (
	"github.com/sportgraphs/rabbit"
	"encoding/json"
	"github.com/streadway/amqp"
)

type MiddlewareHandler func(message Message) error

type Middleware interface {
	Handle(message Message, next MiddlewareHandler) error
}

type DelegatesMessageToHandlerMiddleware struct {
	commandHandlerRegistry CommandHandlerRegistry
}

func NewDelegatesMessageToHandlerMiddleware(registry CommandHandlerRegistry) *DelegatesMessageToHandlerMiddleware {
	return &DelegatesMessageToHandlerMiddleware{
		commandHandlerRegistry: registry,
	}
}

func (middleware *DelegatesMessageToHandlerMiddleware) Handle(message Message, next MiddlewareHandler) error {
	handler, err := middleware.commandHandlerRegistry.Get(Command(message))

	if err != nil {
		return err
	}

	handler(Command(message))

	return next(message)
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

func (middleware *NotifiesMessageSubscribersMiddleware) Handle(message Message, next MiddlewareHandler) error {
	subscribers, err := middleware.eventHandlerRegistry.Get(Event(message))

	if err != nil {
		return next(message)
	}

	for _, subscriber := range subscribers {
		subscriber(Event(message))
	}

	return next(message)
}

func (middleware *AlwaysPublishesMessagesMiddleware) Handle(message Message, next MiddlewareHandler) error {
	err := middleware.publish(message)

	if err != nil {
		return err
	}

	return next(message)
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

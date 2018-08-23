package godatabus

import (
	"github.com/streadway/amqp"
	"github.com/sportgraphs/rabbit"
	"encoding/json"
	"log"
	"reflect"
	"context"
)

type Envelope struct {
	Type    string          `json:"type"`
	Payload string          `json:"payload"`
}

func InitializeRabbitMessageHandler(queueName string, messageBus MessageBus, mq rabbit.MQ, resolver NameResolver, ctx context.Context, logger *log.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}

	mq.SetConsumerHandler(queueName, ctx, func(ctx context.Context, delivery amqp.Delivery) {
		var envelope Envelope

		if err := json.Unmarshal(delivery.Body, &envelope); err != nil {
			delivery.Reject(true)

			logger.Println("unable to process raw message")

			return
		}

		messageType, err := resolver.Resolve(envelope.Type)
		if err != nil {
			delivery.Reject(true)

			logger.Printf("unable to resolve message type %s", envelope.Type)

			return
		}
		messageValue := reflect.New(messageType)
		message := messageValue.Interface()

		if err := json.Unmarshal([]byte(envelope.Payload), &message); err != nil {
			delivery.Reject(true)

			logger.Printf("unable to process message payload for %s (%s)", envelope.Type, err)

			return
		}

		_, err = messageBus.Handle(ctx, Message(message.(Message)))
		if err != nil {
			delivery.Reject(true)

			logger.Printf("unable to handle message for %s (%s)", envelope.Type, err)

			return
		}

		if err = delivery.Ack(false); err != nil {
			return
		}
	})
}

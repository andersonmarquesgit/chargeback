package consumers

import (
	"batch/internal/infrastructure/logging"
	"batch/internal/infrastructure/rabbitmq"
	"batch/internal/interfaces/events/handlers"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	Connection  *amqp.Connection
	Exchange    string
	QueueSuffix string
	RoutingKeys []string
}

type Consumers struct {
	ChargebackBatchEventConsumer *Consumer
	// Add more consumers here
}

func NewConsumers(conn *amqp.Connection, eventHandlers *handlers.EventHandlers) (*Consumers, error) {
	chargebackBatchEventConsumer, err := NewChargebackBatchEventConsumer(conn)

	if err != nil {
		log.Fatalf("Could not execute batch file consumers: %v", err)
	}

	// Watch the queue and consumers events
	go func() {
		err = chargebackBatchEventConsumer.Listen(eventHandlers.ChargebackBatchEventHandler.HandleChargebackOpenedEvent)

		logging.Infof("Listening to batch file consumers: %v", err)

		if err != nil {
			log.Fatalf("Error listening to batch file consumers: %v", err)
		}
	}()

	return &Consumers{
		ChargebackBatchEventConsumer: chargebackBatchEventConsumer,
	}, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.Connection.Channel()
	if err != nil {
		return err
	}

	err = rabbitmq.DeclareExchange(channel, consumer.Exchange, "topic")
	if err != nil {
		return err
	}

	for _, key := range consumer.RoutingKeys {
		queueName := fmt.Sprintf("%s-%s", key, consumer.QueueSuffix)
		_, err := rabbitmq.DeclareQueue(channel, queueName)
		if err != nil {
			return err
		}
		err = channel.QueueBind(queueName, key, consumer.Exchange, false, nil)
		if err != nil {
			return err
		}
		log.Printf("Queue %s bound to key %s in exchange %s", queueName, key, consumer.Exchange)
	}

	return nil
}

func (consumer *Consumer) Listen(handler func(amqp.Delivery) error) error {
	ch, err := consumer.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	for _, key := range consumer.RoutingKeys {
		queueName := fmt.Sprintf("%s-%s", key, consumer.QueueSuffix)

		// Bind the queue to the exchange with the routing key
		err = ch.QueueBind(queueName, key, consumer.Exchange, false, nil)
		if err != nil {
			return err
		}

		// Start consuming messages
		msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
		if err != nil {
			return err
		}

		go func() {
			for msg := range msgs {
				addLoggerHooks(msg)

				// Handle the message with the handler
				if err := handler(msg); err != nil {
					logging.Infof("could not handle message: %v", err)
				}
			}
		}()
	}

	select {} // This blocks forever, similar to <-forever in the other example.

	return nil
}

func addLoggerHooks(msg amqp.Delivery) {
	// Extract the headers from the message
	traceID, traceIDExists := msg.Headers["request-trace-id"]

	// Adicionar traceID ao logger
	if traceIDExists {
		logging.Logger.AddHook(&logging.TraceIDHook{
			TraceIDKey: "request-trace-id",
			Context:    context.WithValue(context.Background(), "request-trace-id", traceID),
		})
	}
}

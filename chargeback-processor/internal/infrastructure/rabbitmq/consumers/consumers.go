package consumers

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"processor/internal/infrastructure/logging"
	"processor/internal/infrastructure/rabbitmq"
	"processor/internal/interfaces/events/handlers"
)

type Consumer struct {
	Connection  *amqp.Connection
	Exchange    string
	QueueSuffix string
	RoutingKeys []string
}

type Consumers struct {
	ChargebackOpenedEventConsumer *Consumer
	// Add more consumers here
}

func NewConsumers(conn *amqp.Connection, eventHandlers *handlers.EventHandlers) (*Consumers, error) {
	chargebackOpenedEventConsumer, err := NewChargebackOpenedEventConsumer(conn)

	if err != nil {
		log.Fatalf("Could not chargeback opened consumers: %v", err)
	}

	// Watch the queue and consumers events
	go func() {
		err = chargebackOpenedEventConsumer.Listen(eventHandlers.ChargebackOpenedEventHandler.HandleChargebackOpenedEvent)

		logging.Infof("Listening to chargeback opened consumers: %v", err)

		if err != nil {
			log.Fatalf("Error listening to chargeback opened consumers: %v", err)
		}
	}()

	return &Consumers{
		ChargebackOpenedEventConsumer: chargebackOpenedEventConsumer,
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
	country, countryExists := msg.Headers["country"]

	// Adicionar traceID e country ao logger
	if traceIDExists {
		logging.Logger.AddHook(&logging.TraceIDHook{
			TraceIDKey: "request-trace-id",
			Context:    context.WithValue(context.Background(), "request-trace-id", traceID),
		})
	}
	if countryExists {
		logging.Logger.AddHook(&logging.CountryHook{
			CountryKey: "country",
			Context:    context.WithValue(context.Background(), "country", country),
		})
	}
}

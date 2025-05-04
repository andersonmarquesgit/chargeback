package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareExchange(ch *amqp.Channel, exchangeName, exchangeType string) error {
	return ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable (messages survive broker restart): auto-ack
		false,        // auto-deleted (delete when no consumers): exclusive
		false,        // internal (used by other exchanges): no-wait
		false,        // no-wait (do not wait for the response): arguments
		nil,          // arguments (optional)
	)
}

func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		false,     // durable?
		false,     // delete when unused?
		true,      // exclusive?
		false,     // no-wait?
		nil,       // arguments?
	)
}

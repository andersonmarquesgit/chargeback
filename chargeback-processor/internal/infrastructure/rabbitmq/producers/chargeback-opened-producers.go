package producers

import (
	"api/internal/infrastructure/logging"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewChargebackOpenedProducer(conn *amqp.Connection) (*Producer, error) {
	producer := Producer{
		Connection:  conn,
		Exchange:    "chargeback-opened-exchange",
		QueueSuffix: "chargeback.opened",
	}

	err := producer.setup("topic")
	if err != nil {
		return nil, err
	}

	return &producer, nil
}

func (p *Producer) PublishChargebackOpenedEvent(message []byte, userID, transactionID, traceID string) error {
	channel, err := p.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	routingKey := ""
	logging.Infof("Publishing chargeback opened event to topic: %s", p.Exchange)

	err = channel.Publish(
		p.Exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         message,
			Headers: amqp.Table{
				"event":            "chargeback-opened",
				"user-id":          userID,
				"transaction-id":   transactionID,
				"request-trace-id": traceID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish chargeback opened event: %w", err)
	}

	return nil
}

package producers

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"processor/internal/infrastructure/logging"
)

func NewChargebackBatchProducer(conn *amqp.Connection) (*Producer, error) {
	producer := Producer{
		Connection:  conn,
		Exchange:    "chargeback-batch-exchange",
		QueueSuffix: "chargeback.batch",
	}

	err := producer.setup("topic")
	if err != nil {
		return nil, err
	}

	return &producer, nil
}

func (p *Producer) PublishChargebackBatchEvent(message []byte, traceID string) error {
	channel, err := p.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	routingKey := ""
	logging.Infof("Publishing chargeback batch event to topic: %s", p.Exchange)

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
				"event":            "chargeback-batch",
				"request-trace-id": traceID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish chargeback batch event: %w", err)
	}

	return nil
}

package consumers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

func NewChargebackBatchEventConsumer(conn *amqp.Connection) (*Consumer, error) {
	routingKeys := strings.Split("", ",")

	consumer := Consumer{
		Connection:  conn,
		Exchange:    "chargeback-batch-exchange",
		QueueSuffix: "chargeback-batch",
		RoutingKeys: routingKeys,
	}

	err := consumer.setup()
	if err != nil {
		return &Consumer{}, err
	}

	return &consumer, nil
}

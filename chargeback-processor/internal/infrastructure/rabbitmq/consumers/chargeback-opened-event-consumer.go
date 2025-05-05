package consumers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

func NewChargebackOpenedEventConsumer(conn *amqp.Connection) (*Consumer, error) {
	routingKeys := strings.Split("", ",")

	consumer := Consumer{
		Connection:  conn,
		Exchange:    "chargeback-opened-exchange",
		QueueSuffix: "chargeback-opened",
		RoutingKeys: routingKeys,
	}

	err := consumer.setup()
	if err != nil {
		return &Consumer{}, err
	}

	return &consumer, nil
}

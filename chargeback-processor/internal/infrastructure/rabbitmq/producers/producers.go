package producers

import (
	"api/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Connection  *amqp.Connection
	Exchange    string
	QueueSuffix string
}

type Producers struct {
	ChargebackOpenedProducer *Producer
	// Add more producers here
}

func NewProducers(conn *amqp.Connection) (*Producers, error) {
	chargebackOpenedProducer, err := NewChargebackOpenedProducer(conn)
	if err != nil {
		return nil, err
	}

	return &Producers{
		ChargebackOpenedProducer: chargebackOpenedProducer,
	}, nil
}

func (producer *Producer) setup(exchangeType string) error {
	channel, err := producer.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return rabbitmq.DeclareExchange(channel, producer.Exchange, exchangeType)
}

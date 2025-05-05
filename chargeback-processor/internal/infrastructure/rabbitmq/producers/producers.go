package producers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"processor/internal/infrastructure/rabbitmq"
)

type Producer struct {
	Connection  *amqp.Connection
	Exchange    string
	QueueSuffix string
}

type Producers struct {
	ChargebackBatchProducer *Producer
	// Add more producers here
}

func NewProducers(conn *amqp.Connection) (*Producers, error) {
	chargebackBatchProducer, err := NewChargebackBatchProducer(conn)
	if err != nil {
		return nil, err
	}

	return &Producers{
		ChargebackBatchProducer: chargebackBatchProducer,
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

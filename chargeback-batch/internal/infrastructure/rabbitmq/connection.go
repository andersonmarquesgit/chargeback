package rabbitmq

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const MAX_RETRY_RABBITMQ_CONNECTION = 5

func Connect(url string) (*amqp.Connection, error) {
	var retryConnection int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(url)
		if err != nil {
			log.Println("RabbitMQ not yet ready ...")
			retryConnection++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if retryConnection > MAX_RETRY_RABBITMQ_CONNECTION {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(retryConnection), 2) * float64(time.Second))
		log.Println("backing off ...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

package event

import (
	"inventory/internal/config"
	"log"

	"github.com/streadway/amqp"
)

const (
	ExchangeName = "monorepo_exchange"
)

// Publisher sends messages to the exchange.
func Publisher(ch *amqp.Channel, routingKey string, body []byte) error {
	err := ch.Publish(
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func New(cfg config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.AMQPHost)
	if err != nil {
		log.Printf("%s\n", err)
		return nil, err
	}

	return conn, nil
}

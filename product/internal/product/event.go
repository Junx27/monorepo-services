package product

import (
	"log"
	"product/pkg/event"

	"github.com/streadway/amqp"
)

type ProductEvent struct {
	Channel *amqp.Channel
}

func NewProductEvent(ch *amqp.Channel) ProductEvent {
	return ProductEvent{
		Channel: ch,
	}
}

func (pe *ProductEvent) SubscribeSomething() {
	q, err := pe.Channel.QueueDeclare(
		"something", // random queue name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = pe.Channel.QueueBind(
		q.Name,             // queue name
		"monorepo",         // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := pe.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for msg := range msgs {
		pe.handleConsumeSomething(msg)
	}
}

func (pe *ProductEvent) handleConsumeSomething(msg amqp.Delivery) {
	// do something

}

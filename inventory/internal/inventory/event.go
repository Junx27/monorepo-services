package inventory

import (
	"fmt"
	"inventory/pkg/event"
	"log"

	"github.com/streadway/amqp"
)

type InventoryEvent struct {
	Channel *amqp.Channel
}

func NewInventoryEvent(ch *amqp.Channel) InventoryEvent {
	return InventoryEvent{
		Channel: ch,
	}
}

func (ie *InventoryEvent) SubscribeSomething() {
	q, err := ie.Channel.QueueDeclare(
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

	err = ie.Channel.QueueBind(
		q.Name,             // queue name
		"monorepo",         // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := ie.Channel.Consume(
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
		fmt.Println(msg)
	}
}

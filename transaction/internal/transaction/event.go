package transaction

import (
	"fmt"
	"log"
	"transaction/pkg/event"

	"github.com/streadway/amqp"
)

type TransactionEvent struct {
	Channel *amqp.Channel
}

func NewTransactionEvent(ch *amqp.Channel) TransactionEvent {
	return TransactionEvent{
		Channel: ch,
	}
}

func (te *TransactionEvent) SubscribeSomething() {
	q, err := te.Channel.QueueDeclare(
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

	err = te.Channel.QueueBind(
		q.Name,             // queue name
		"reduce.product",   // routing key
		event.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	msgs, err := te.Channel.Consume(
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
		te.handleConsumeSomething(msg)
	}
}

func (te *TransactionEvent) handleConsumeSomething(msg amqp.Delivery) {
	message := string(msg.Body)

	fmt.Printf("Received message from transaction service: %s\n", message)
}

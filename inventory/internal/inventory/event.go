package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inventory/pkg/event"
	"log"
	"net/http"

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

func (ie *InventoryEvent) SubscribeTransactionPaid() {
	q, err := ie.Channel.QueueDeclare(
		"update_transaction_paid", // random queue name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = ie.Channel.QueueBind(
		q.Name,                    // queue name
		"update.transaction.paid", // routing key
		event.ExchangeName,        // exchange
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
		ie.handleConsumeTransactionPaid(msg)
	}
}
func (ie *InventoryEvent) SubscribeProductIncreased() {
	q, err := ie.Channel.QueueDeclare(
		"update_product_increased", // random queue name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = ie.Channel.QueueBind(
		q.Name,                     // queue name
		"update.product.increased", // routing key
		event.ExchangeName,         // exchange
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
		ie.handleConsumeProductIncreased(msg)
	}
}

func (ie *InventoryEvent) handleConsumeTransactionPaid(msg amqp.Delivery) {
	var inventoryTransaction InventoryTransaction
	err := json.Unmarshal(msg.Body, &inventoryTransaction)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	fmt.Printf("Received message from inventory service: %s\n", string(msg.Body))
	fmt.Printf("Inventory Transaction Data: %+v\n", inventoryTransaction.ProductID)

	ie.createInventory(inventoryTransaction)
}
func (ie *InventoryEvent) handleConsumeProductIncreased(msg amqp.Delivery) {
	message := string(msg.Body)

	fmt.Printf("Received message from inventory service: %s\n", message)
	fmt.Println("Update inventory success")
	fmt.Println("Successfully order to cancel")
}

func (ie *InventoryEvent) createInventory(inventoryTransaction InventoryTransaction) {
	url := fmt.Sprintf("http://localhost:8001/inventory")

	fmt.Printf("Sending POST request to: %s\n", url)

	postBody, err := json.Marshal(inventoryTransaction)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("Error sending HTTP POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Product stock successfully reduce!")
	} else {
		log.Printf("Failed to reduce product stock: Status Code %d, %s", resp.StatusCode, resp.Status)
	}
}

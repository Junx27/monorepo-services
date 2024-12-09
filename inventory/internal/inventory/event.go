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

func (ie *InventoryEvent) handleConsumeTransactionPaid(msg amqp.Delivery) {
	var inventoryTransaction InventoryTransaction
	err := json.Unmarshal(msg.Body, &inventoryTransaction)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	fmt.Printf("Received message from inventory service: %s\n", inventoryTransaction)
	ie.createInventory(inventoryTransaction)

}

func (ie *InventoryEvent) createInventory(inventoryTransaction InventoryTransaction) {
	// Membentuk URL dengan product_id
	url := fmt.Sprintf("http://localhost:8001/inventory")

	// Menampilkan URL yang digunakan untuk debug
	fmt.Printf("Sending POST request to: %s\n", url)

	// Menyiapkan body request
	postBody, err := json.Marshal(inventoryTransaction)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return
	}

	// Mengirimkan POST request ke endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("Error sending HTTP POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Menampilkan status kode response
	if resp.StatusCode == 200 {
		fmt.Println("Product stock successfully increase!")
	} else {
		log.Printf("Failed to reduce product stock: Status Code %d, %s", resp.StatusCode, resp.Status)
	}
}

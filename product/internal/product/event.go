package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product/pkg/event" // pastikan ini mengarah ke paket yang benar

	"github.com/streadway/amqp"
)

// Struktur untuk pesan yang dikirim ke RabbitMQ
type ProductTransaction struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ProductEvent struct {
	Channel *amqp.Channel
}

func NewProductEvent(ch *amqp.Channel) ProductEvent {
	return ProductEvent{
		Channel: ch,
	}
}

func (pe *ProductEvent) SubscribeReduceStock() {
	q, err := pe.Channel.QueueDeclare(
		"create_transaction_success", // nama antrian
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring queue: %s", err)
	}

	err = pe.Channel.QueueBind(
		q.Name,                       // queue name
		"create.transaction.success", // routing key
		event.ExchangeName,           // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error binding queue: %s", err)
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
		log.Fatalf("Error consuming messages: %s", err)
	}

	// Mengkonsumsi pesan yang diterima dari RabbitMQ
	for msg := range msgs {
		pe.handleConsumeReduceStock(msg)
	}
}
func (pe *ProductEvent) SubscribeIncreaseStock() {
	q, err := pe.Channel.QueueDeclare(
		"update_transaction_cancel", // nama antrian
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring queue: %s", err)
	}

	err = pe.Channel.QueueBind(
		q.Name,                      // queue name
		"update.transaction.cancel", // routing key
		event.ExchangeName,          // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error binding queue: %s", err)
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
		log.Fatalf("Error consuming messages: %s", err)
	}

	// Mengkonsumsi pesan yang diterima dari RabbitMQ
	for msg := range msgs {
		pe.handleConsumeIncreaseStock(msg)
	}
}

func (pe *ProductEvent) handleConsumeReduceStock(msg amqp.Delivery) {
	// Parsing pesan JSON dari RabbitMQ
	var productTransaction ProductTransaction
	err := json.Unmarshal(msg.Body, &productTransaction)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	// Mengirim HTTP POST request untuk memproses transaksi produk

	pe.reduceStockProduct(productTransaction)
}

func (pe *ProductEvent) handleConsumeIncreaseStock(msg amqp.Delivery) {
	// Parsing pesan JSON dari RabbitMQ
	var productTransaction ProductTransaction
	err := json.Unmarshal(msg.Body, &productTransaction)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	// Mengirim HTTP POST request untuk memproses transaksi produk

	fmt.Printf("Received message from product service: %s\n", productTransaction)
	pe.increaseStockProduct(productTransaction)
}

// Fungsi untuk mengirimkan HTTP PUT request untuk memperbarui stok produk
func (pe *ProductEvent) increaseStockProduct(productTransaction ProductTransaction) {
	// Membentuk URL dengan product_id
	url := fmt.Sprintf("http://localhost:8000/product/%s/increaseStock", productTransaction.ProductID)

	// Menampilkan URL yang digunakan untuk debug
	fmt.Printf("Sending POST request to: %s\n", url)

	// Menyiapkan body request
	postBody, err := json.Marshal(productTransaction)
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

// Fungsi untuk mengirimkan HTTP POST request untuk mengurangi stok produk
func (pe *ProductEvent) reduceStockProduct(productTransaction ProductTransaction) {
	// Membentuk URL dengan product_id
	url := fmt.Sprintf("http://localhost:8000/product/%s/reduceStock", productTransaction.ProductID)

	// Menampilkan URL yang digunakan untuk debug
	fmt.Printf("Sending POST request to: %s\n", url)

	// Menyiapkan body request
	postBody, err := json.Marshal(productTransaction)
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
		fmt.Println("Product stock successfully reduced!")
	} else {
		log.Printf("Failed to reduce product stock: Status Code %d, %s", resp.StatusCode, resp.Status)
	}
}

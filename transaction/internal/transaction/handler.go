package transaction

import (
	"encoding/json"
	"net/http"
	"transaction/internal/config"
	"transaction/pkg/event"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/streadway/amqp"
)

type Handler struct {
	cfg config.Config
	ch  *amqp.Channel
}

func NewHandler(cfg config.Config, ch *amqp.Channel) Handler {
	return Handler{
		cfg: cfg,
		ch:  ch,
	}
}

type Product struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
}

// GetTransaction retrieves a transaction from the database by ID
func (h *Handler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	transaction, err := h.GetTransactionFromDB(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction"})
		}
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// CreateTransaction creates a new transaction
func (h *Handler) CreateTransaction(c *gin.Context) {
	var transaction Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	err := StoreTransaction(c.Request.Context(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	product := Product{
		UserID:    transaction.UserID,
		ProductID: transaction.ProductID,
	}

	// Convert product to JSON
	productData, err := json.Marshal(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unmarshal error"})
		return
	}

	// Send event about the successful creation of the transaction
	err = event.Publisher(h.ch, "create.transaction.success", productData)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
}

// UpdateTransactionStatus is a generic function to update the status of a transaction

func (h *Handler) UpdateTransaction(c *gin.Context, status string) {
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Step 1: Get transaction details from the database
	transaction, err := h.GetTransactionFromDB(c, transactionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction"})
		}
		return
	}

	// Step 2: Update the status of the transaction in the database
	err = UpdateTransactionStatus(c.Request.Context(), transactionID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction status"})
		return
	}

	// Step 3: Prepare the event data
	eventData := struct {
		TransactionID string `json:"transaction_id"`
		UserID        string `json:"user_id"`
		ProductID     string `json:"product_id"`
		Quantity      int    `json:"quantity"`
	}{
		TransactionID: transaction.UUID,
		UserID:        transaction.UserID,
		ProductID:     transaction.ProductID,
		Quantity:      1,
	}

	// Convert event data to JSON
	eventBytes, err := json.Marshal(eventData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal event data"})
		return
	}

	// Step 4: Publish the event
	var eventName string
	if status == "paid" {
		eventName = "update.transaction.paid"
	} else if status == "cancel" {
		eventName = "update.transaction.cancel"
	}

	err = event.Publisher(h.ch, eventName, eventBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to publish event"})
		return
	}

	// Step 5: Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction status updated to " + status,
	})
}

// UpdateTransactionPaid updates the transaction status to "paid"
func (h *Handler) UpdateTransactionPaid(c *gin.Context) {
	h.UpdateTransaction(c, "paid")
}

// UpdateTransactionCancel updates the transaction status to "cancel"
func (h *Handler) UpdateTransactionCancel(c *gin.Context) {
	h.UpdateTransaction(c, "cancel")
}

package transaction

import (
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

func (h *Handler) CreateTransaction(c *gin.Context) {
	var transaction Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	err := StoreTransaction(c.Request.Context(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	event.Publisher(h.ch, "create.transaction", []byte("transaction successfully"))
	c.JSON(http.StatusOK, gin.H{
		"message": "success creating transaction",
	})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	// update transaction with status: PAID
	// call service product to reduce the stock by 1
	// add new item to the user inventory by calling inventory service
	event.Publisher(h.ch, "update.transaction", []byte("transaction update"))
	// api driven atau event driven
	c.JSON(http.StatusOK, gin.H{
		"message": "success update transaction",
	})
}

package product

import (
	"net/http"
	"product/internal/config"
	"product/pkg/event"

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

func (h *Handler) GetProductHandler(c *gin.Context) {
	productID := c.Param("id")
	product, err := GetProduct(c.Request.Context(), productID)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	event.Publisher(h.ch, "update.product", []byte("product update"))
	// api driven atau event driven
	c.JSON(http.StatusOK, gin.H{
		"message": "success update product",
	})

}

func (h *Handler) ReduceStock(c *gin.Context) {
	event.Publisher(h.ch, "reduce.product", []byte("product reduce"))
	// api driven atau event driven
	c.JSON(http.StatusOK, gin.H{
		"message": "success reduce product",
	})

}

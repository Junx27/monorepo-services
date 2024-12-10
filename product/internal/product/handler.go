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
	var req Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	err := UpdateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated successfully",
	})
}

func (h *Handler) ReduceStock(c *gin.Context) {
	productID := c.Param("id")

	err := ReduceStock(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to reduce stock",
		})
		return
	}
	event.Publisher(h.ch, "transaction.success", []byte("update inventory success"))

	c.JSON(http.StatusOK, gin.H{
		"message": "product stock reduced by 1",
	})
}

func (h *Handler) IncreaseStock(c *gin.Context) {
	productID := c.Param("id")

	err := IncreaseStock(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to increase stock",
		})
		return
	}

	event.Publisher(h.ch, "update.product.increased", []byte("update inventory success"))

	c.JSON(http.StatusOK, gin.H{
		"message": "product stock increased by 1",
	})
}

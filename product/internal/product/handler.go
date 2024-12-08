package product

import (
	"net/http"
	"product/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) Handler {
	return Handler{
		cfg: cfg,
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

}

func (h *Handler) ReduceStock(c *gin.Context) {

}

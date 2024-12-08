package inventory

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) GetInventory(c *gin.Context) {

}

func (h *Handler) CreateInventory(c *gin.Context) {

}

func (h *Handler) UpdateInventory(c *gin.Context) {

}

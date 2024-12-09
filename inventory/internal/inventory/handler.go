package inventory

import (
	"inventory/internal/config"
	"log"
	"net/http"

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

func (h *Handler) GetInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	inventory, err := GetInventory(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Inventory"})
		}
		return
	}

	c.JSON(http.StatusOK, inventory)

}

func (h *Handler) CreateInventory(c *gin.Context) {
	var inventoryRequest Inventory

	// Bind the JSON body to the inventory struct
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
		})
		return
	}

	// Store the inventory in the database (with a quantity of 1)
	err := StoreInventory(c.Request.Context(), inventoryRequest)
	if err != nil {
		log.Printf("Error storing inventory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully created inventory with quantity 1",
	})
}

func (h *Handler) UpdateInventory(c *gin.Context) {
	// Get the inventory ID from the URL parameters
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Update the inventory in the database by adding 1 to the quantity
	err := UpdateInventory(c.Request.Context(), id)
	if err != nil {
		// Return a 404 error if the inventory ID does not exist
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		} else {
			// Return an internal server error for other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		}
		return
	}

	// Return a success response after updating the inventory
	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory updated successfully",
	})
}

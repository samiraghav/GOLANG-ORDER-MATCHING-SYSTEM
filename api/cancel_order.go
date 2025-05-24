package api

import (
	"net/http"
	"order-matching-engine/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CancelOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := db.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.Status != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order is already filled or canceled"})
		return
	}

	err = db.UpdateOrderStatus(orderID, "canceled", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order canceled"})
}

package api

import (
	"fmt"
	"net/http"
	"order-matching-engine/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}
	fmt.Println("orderID", orderID)

	order, err := db.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found", "details": err.Error()})
		return
	}

	fmt.Println("order", order)

	c.JSON(http.StatusOK, gin.H{
		"order_id":           order.ID,
		"symbol":             order.Symbol,
		"side":               order.Side,
		"type":               order.Type,
		"price":              order.Price,
		"quantity":           order.Quantity,
		"remaining_quantity": order.RemainingQty,
		"status":             order.Status,
		"created_at":         order.CreatedAt,
	})
}

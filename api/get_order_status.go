package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := db.GetOrderByID(orderID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Order not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Order status", map[string]interface{}{
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

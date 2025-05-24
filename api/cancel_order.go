package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CancelOrder(c *gin.Context) {
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

	if order.Status != "open" {
		utils.SendError(c, http.StatusBadRequest, "Order is already filled or canceled")
		return
	}

	if err := db.UpdateOrderStatus(orderID, "canceled", 0); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Order canceled", nil)
}

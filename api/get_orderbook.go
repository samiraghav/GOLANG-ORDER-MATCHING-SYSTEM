package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/utils"

	"github.com/gin-gonic/gin"
)

func GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		utils.SendError(c, http.StatusBadRequest, "Symbol is required")
		return
	}

	buyOrders, err := db.GetOrdersBySymbolAndSide(symbol, "buy")
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch buy orders")
		return
	}

	sellOrders, err := db.GetOrdersBySymbolAndSide(symbol, "sell")
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch sell orders")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Order book fetched", map[string]interface{}{
		"symbol":      symbol,
		"buy_orders":  buyOrders,
		"sell_orders": sellOrders,
	})
}

package api

import (
	"net/http"
	"order-matching-engine/db"

	"github.com/gin-gonic/gin"
)

func GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	buyOrders, err := db.GetOrdersBySymbolAndSide(symbol, "buy")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch buy orders"})
		return
	}

	sellOrders, err := db.GetOrdersBySymbolAndSide(symbol, "sell")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sell orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol":      symbol,
		"buy_orders":  buyOrders,
		"sell_orders": sellOrders,
	})
}

package api

import (
	"net/http"
	"order-matching-engine/db"

	"github.com/gin-gonic/gin"
)

func ListTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	trades, err := db.GetTradesBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol": symbol,
		"trades": trades,
	})
}

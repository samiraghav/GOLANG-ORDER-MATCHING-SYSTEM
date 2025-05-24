package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/utils"

	"github.com/gin-gonic/gin"
)

func ListTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		utils.SendError(c, http.StatusBadRequest, "Symbol is required")
		return
	}

	trades, err := db.GetTradesBySymbol(symbol)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch trades")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Trades fetched", map[string]interface{}{
		"symbol": symbol,
		"trades": trades,
	})
}

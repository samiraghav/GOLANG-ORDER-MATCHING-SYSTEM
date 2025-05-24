package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/models"
	"order-matching-engine/service"
	"order-matching-engine/utils"

	"github.com/gin-gonic/gin"
)

type PlaceOrderRequest struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Side     string  `json:"side" binding:"required,oneof=buy sell"`
	Type     string  `json:"type" binding:"required,oneof=limit market"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}

func PlaceOrder(c *gin.Context) {
	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	order := models.Order{
		Symbol:       req.Symbol,
		Side:         req.Side,
		Type:         req.Type,
		Price:        req.Price,
		Quantity:     req.Quantity,
		RemainingQty: req.Quantity,
		Status:       "open",
	}

	insertedID, err := db.InsertOrder(&order)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to place order")
		return
	}

	order.ID = insertedID
	go service.MatchOrder(order)

	utils.SendSuccess(c, http.StatusOK, "Order placed", map[string]interface{}{
		"order_id": insertedID,
	})
}

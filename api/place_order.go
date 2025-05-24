package api

import (
	"net/http"
	"order-matching-engine/db"
	"order-matching-engine/models"
	"order-matching-engine/service"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to place order"})
		return
	}

	order.ID = insertedID
	go service.MatchOrder(order)

	c.JSON(http.StatusOK, gin.H{"message": "order placed", "order_id": insertedID})
}

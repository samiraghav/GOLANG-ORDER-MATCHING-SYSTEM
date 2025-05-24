package router

import (
	"order-matching-engine/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/orders", api.PlaceOrder)
	r.DELETE("/orders/:id", api.CancelOrder)
	r.GET("/orderbook", api.GetOrderBook)
	r.GET("/trades", api.ListTrades)
	r.GET("/orders/:id", api.GetOrderStatus)

	return r
}

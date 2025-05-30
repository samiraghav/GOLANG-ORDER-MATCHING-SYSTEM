package service

import (
	"order-matching-engine/db"
	"order-matching-engine/models"
	"order-matching-engine/utils"
	"sort"
	"time"
)

func MatchOrder(order models.Order) {
	// Get opposite side orders from the order book
	oppositeSide := "buy"
	if order.Side == "buy" {
		oppositeSide = "sell"
	}

	orders, err := db.GetOpenOrders(order.Symbol, oppositeSide)
	if err != nil {
		utils.LogError("MatchOrder:GetOpenOrders", map[string]interface{}{
			"error":  err.Error(),
			"symbol": order.Symbol,
			"side":   oppositeSide,
		})
		return
	}

	if len(orders) == 0 {
		if order.Type == "market" {
			_ = db.UpdateOrderStatus(order.ID, "canceled", order.RemainingQty)
		}
		return
	}

	// Sort by price-time priority
	sort.SliceStable(orders, func(i, j int) bool {
		const layout = "2006-01-02 15:04:05"
		timeI, err1 := time.Parse(layout, orders[i].CreatedAt)
		timeJ, err2 := time.Parse(layout, orders[j].CreatedAt)

		if err1 != nil || err2 != nil {
			utils.LogError("MatchOrder:TimeParse", map[string]interface{}{
				"err1": err1.Error(),
				"err2": err2.Error(),
			})
			return false
		}

		if orders[i].Price == orders[j].Price {
			return timeI.Before(timeJ)
		}

		if order.Side == "buy" {
			// higher bid first
			return orders[i].Price > orders[j].Price
		}
		// lower ask first
		return orders[i].Price < orders[j].Price
	})

	remainingQty := order.RemainingQty

	for _, bookOrder := range orders {
		if remainingQty == 0 {
			break
		}

		// Check price match condition
		priceMatch := false
		if order.Type == "market" {
			priceMatch = true
		} else if order.Side == "buy" && order.Price >= bookOrder.Price {
			priceMatch = true
		} else if order.Side == "sell" && order.Price <= bookOrder.Price {
			priceMatch = true
		}

		if !priceMatch {
			continue
		}

		// Calculate match quantity
		matchQty := min(remainingQty, bookOrder.RemainingQty)
		tradePrice := bookOrder.Price
		if order.Type == "market" {
			// market order takes book price
			tradePrice = bookOrder.Price
		}

		// Record trade
		trade := models.Trade{
			BuyOrderID:  ifBuy(order, bookOrder),
			SellOrderID: ifSell(order, bookOrder),
			Symbol:      order.Symbol,
			Price:       tradePrice,
			Quantity:    matchQty,
		}

		if err := db.InsertTrade(trade); err != nil {
			utils.LogError("MatchOrder:InsertTrade", map[string]interface{}{
				"error":  err.Error(),
				"symbol": order.Symbol,
			})
			continue
		}

		newRemaining := bookOrder.RemainingQty - matchQty

		if err := db.DecreaseOrderQty(bookOrder.ID, matchQty); err != nil {
			utils.LogError("MatchOrder:DecreaseOrderQty", map[string]interface{}{
				"error":     err.Error(),
				"order_id":  bookOrder.ID,
				"matched":   matchQty,
				"remaining": newRemaining,
			})
		}

		if newRemaining == 0 {
			if err := db.UpdateOrderStatus(bookOrder.ID, "filled", 0); err != nil {
				utils.LogError("MatchOrder:FillBookOrder", map[string]interface{}{
					"error":    err.Error(),
					"order_id": bookOrder.ID,
				})
			}
		}

		// Update incoming order state
		remainingQty -= matchQty
	}

	// Final status for incoming order
	if remainingQty == 0 {
		_ = db.UpdateOrderStatus(order.ID, "filled", 0)
	} else if order.Type == "limit" {
		_ = db.UpdateOrderStatus(order.ID, "open", remainingQty)
	} else {
		// market order partial fill
		_ = db.UpdateOrderStatus(order.ID, "canceled", 0)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ifBuy(o1, o2 models.Order) int64 {
	if o1.Side == "buy" {
		return o1.ID
	}
	return o2.ID
}

func ifSell(o1, o2 models.Order) int64 {
	if o1.Side == "sell" {
		return o1.ID
	}
	return o2.ID
}

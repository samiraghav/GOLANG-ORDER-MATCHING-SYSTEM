package service

import (
	"log"
	"order-matching-engine/db"
	"order-matching-engine/models"
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
	if err != nil || len(orders) == 0 {
		return
	}

	// Sort by price-time priority
	sort.SliceStable(orders, func(i, j int) bool {
		const layout = "2006-01-02 15:04:05"
		timeI, err1 := time.Parse(layout, orders[i].CreatedAt)
		timeJ, err2 := time.Parse(layout, orders[j].CreatedAt)

		if err1 != nil || err2 != nil {
			log.Println("Time parse error in order sorting:", err1, err2)
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
			break
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
		err := db.InsertTrade(trade)
		if err != nil {
			log.Println("Trade insert failed:", err)
			continue
		}

		//  Update book order
		_ = db.DecreaseOrderQty(bookOrder.ID, matchQty)

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

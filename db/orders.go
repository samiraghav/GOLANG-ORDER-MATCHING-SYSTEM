package db

import (
	"order-matching-engine/models"
	"time"
)

func InsertOrder(order *models.Order) (int64, error) {
	query := `INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, order.Symbol, order.Side, order.Type, order.Price, order.Quantity, order.RemainingQty, order.Status, time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetOpenOrders(symbol, side string) ([]models.Order, error) {
	query := `SELECT * FROM orders WHERE symbol = ? AND side = ? AND status = 'open' ORDER BY created_at ASC`
	rows, err := DB.Query(query, symbol, side)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQty, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func InsertTrade(t models.Trade) error {
	query := `INSERT INTO trades (buy_order_id, sell_order_id, symbol, price, quantity, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, t.BuyOrderID, t.SellOrderID, t.Symbol, t.Price, t.Quantity, time.Now())
	return err
}

func DecreaseOrderQty(orderID int64, qty int) error {
	query := `UPDATE orders SET remaining_quantity = remaining_quantity - ? WHERE id = ?`
	_, err := DB.Exec(query, qty, orderID)
	return err
}

func UpdateOrderStatus(orderID int64, status string, remaining int) error {
	query := `UPDATE orders SET status = ?, remaining_quantity = ? WHERE id = ?`
	_, err := DB.Exec(query, status, remaining, orderID)
	return err
}

func GetOrderByID(orderID int64) (*models.Order, error) {
	query := `SELECT * FROM orders WHERE id = ?`
	row := DB.QueryRow(query, orderID)

	var o models.Order
	err := row.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQty, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func GetOrdersBySymbolAndSide(symbol, side string) ([]models.Order, error) {
	query := `SELECT * FROM orders WHERE symbol = ? AND side = ? AND status = 'open' ORDER BY price DESC, created_at ASC`
	if side == "sell" {
		query = `SELECT * FROM orders WHERE symbol = ? AND side = ? AND status = 'open' ORDER BY price ASC, created_at ASC`
	}

	rows, err := DB.Query(query, symbol, side)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQty, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

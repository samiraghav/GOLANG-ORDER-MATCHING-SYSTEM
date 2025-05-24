package db

import (
	"order-matching-engine/models"
)

func GetTradesBySymbol(symbol string) ([]models.Trade, error) {
	query := `SELECT id, buy_order_id, sell_order_id, symbol, price, quantity, created_at 
	          FROM trades 
	          WHERE symbol = ? 
	          ORDER BY created_at DESC 
	          LIMIT 50`

	rows, err := DB.Query(query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var t models.Trade
		err := rows.Scan(&t.ID, &t.BuyOrderID, &t.SellOrderID, &t.Symbol, &t.Price, &t.Quantity, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

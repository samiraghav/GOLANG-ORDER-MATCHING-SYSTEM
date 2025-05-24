package models

type Order struct {
	ID           int64   `json:"id"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	Type         string  `json:"type"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	RemainingQty int     `json:"remaining_quantity"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

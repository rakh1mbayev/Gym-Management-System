package domain

import "time"

type Order struct {
	OrderID    string      `json:"order_id"`
	UserID     int64       `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"update_at"`
}

type OrderItem struct {
	OrderItemID  int64   `json:"order_item_id"`
	OrderID      string  `json:"order_id"`
	ProductID    int64   `json:"product_id"`
	Quantity     int32   `json:"quantity"`
	PricePerItem float32 `json:"price_per_item"`
}

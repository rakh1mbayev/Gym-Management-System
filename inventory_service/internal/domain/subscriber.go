package domain

type OrderCreatedEvent struct {
	OrderID    string  `json:"order_id"`
	UserID     int64   `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	Items      []struct {
		ProductID    int64   `json:"product_id"`
		Quantity     int     `json:"quantity"`
		PricePerItem float32 `json:"price_per_item"`
	} `json:"items"`
}

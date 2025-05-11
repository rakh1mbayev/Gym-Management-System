package domain

import "time"

type Order struct {
	OrderID    string
	UserID     int64
	Items      []OrderItem
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderItem struct {
	OrderItemID  int64
	OrderID      string
	ProductID    int64
	Quantity     int32
	PricePerItem float32
}

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
	OrderItemID  int
	OrderID      string
	ProductID    int
	Quantity     int
	PricePerItem float64
}

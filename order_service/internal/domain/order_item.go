package domain

type OrderItem struct {
	ID           int
	OrderID      int
	ProductID    int
	Quantity     float64
	PricePerItem float64
}

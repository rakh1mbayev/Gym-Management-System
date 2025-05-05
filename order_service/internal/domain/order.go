package domain

type Order struct {
	ID         string
	UserID     string
	Items      []OrderItem
	TotalPrice float64
	Status     string
}

type OrderItem struct {
	ProductID    string
	Quantity     int
	PricePerItem float64
}

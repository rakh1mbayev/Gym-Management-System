package domain

import "time"

type OrderStatus string

const (
	StatusPayed     OrderStatus = "payed"
	StatusDelivered OrderStatus = "delivered"
	StatusCanceled  OrderStatus = "canceled"
)

type Order struct {
	ID        int
	UserID    int
	Name      string
	Price     float64
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	UpdateStatus(id int, status OrderStatus) error
	ListByUser(userID int) ([]Order, error)
}

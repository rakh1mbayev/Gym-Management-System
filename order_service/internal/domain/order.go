package domain

import "time"

type OrderStatus string

const (
	StatusPayed     OrderStatus = "payed"
	StatusDelivered OrderStatus = "delivered"
	StatusCanceled  OrderStatus = "canceled"
)

type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Name      string      `json:"name"`
	Price     float64     `json:"price"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []OrderItem `json:"items"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	UpdateStatus(id int, status OrderStatus) error
	ListByUser(userID int) ([]Order, error)
}

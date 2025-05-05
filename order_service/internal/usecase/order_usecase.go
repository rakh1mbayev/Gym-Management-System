package usecase

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/repository/postgres"
)

type OrderUsecase struct {
	orderRepo postgres.OrderRepository
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []domain.OrderItem, totalPrice float64) (string, error)
	GetOrder(ctx context.Context, orderID string) (*domain.Order, error)
	ListOrders(ctx context.Context, userID string) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}

func NewOrderUsecase(orderRepo postgres.OrderRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem, totalPrice float64) (string, error) {
	// Create an Order struct with all necessary fields
	order := &domain.Order{
		UserID:     userID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     "pending", // Setting the initial status to "pending"
	}

	// Pass the order to the repository layer for saving it to the database
	err := u.orderRepo.Create(order)
	if err != nil {
		return "", err
	}

	return order.ID, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return u.orderRepo.GetByID(orderID)
}

func (u *OrderUsecase) ListOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	return u.orderRepo.ListByUser(userID)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	// Update the order's status
	return u.orderRepo.UpdateStatus(orderID, status)
}

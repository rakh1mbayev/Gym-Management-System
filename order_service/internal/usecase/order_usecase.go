package usecase

import (
	"context"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/proto/inventorypb"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/repository/postgres"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/pkg/nats"
	"log"
)

type OrderUsecase struct {
	orderRepo postgres.OrderRepository
	publisher *nats.NatsPublisher
	inventory inventorypb.InventoryServiceClient
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID int64, items []domain.OrderItem) (string, error)
	GetOrder(ctx context.Context, orderID string) (*domain.Order, error)
	ListOrders(ctx context.Context, userID int64) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}

func NewOrderUsecase(
	orderRepo postgres.OrderRepository,
	publisher *nats.NatsPublisher,
	inventory inventorypb.InventoryServiceClient,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
		publisher: publisher,
		inventory: inventory,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, userID int64, items []domain.OrderItem) (string, error) {
	var totalPrice float64
	var enrichedItems []domain.OrderItem

	for _, item := range items {
		resp, err := u.inventory.GetProduct(ctx, &inventorypb.GetProductRequest{
			ProductId: item.ProductID,
		})
		if err != nil {
			return "", fmt.Errorf("failed to fetch product %d: %w", item.ProductID, err)
		}

		item.PricePerItem = resp.Price
		totalPrice += float64(resp.Price) * float64(item.Quantity)
		enrichedItems = append(enrichedItems, item)
	}

	order := &domain.Order{
		UserID:     userID,
		Items:      enrichedItems,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	if err := u.orderRepo.Create(order); err != nil {
		return "", err
	}

	if err := u.publisher.PublishOrderCreated(order); err != nil {
		log.Printf("Failed to publish order.created: %v", err)
	}

	return order.OrderID, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return u.orderRepo.GetByID(orderID)
}

func (u *OrderUsecase) ListOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	return u.orderRepo.ListByUser(userID)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	// Update the order's status
	return u.orderRepo.UpdateStatus(orderID, status)
}

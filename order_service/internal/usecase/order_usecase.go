package usecase

import (
	"fmt"
	"order_serivce/internal/domain"
)

type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(r domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: r}
}

func (uc *OrderUsecase) Create(order *domain.Order) error {
	// Before saving, check stock for each item
	for _, item := range order.Items {
		resp, err := uc.invClient.CheckStock(ctx, &inventorypb.CheckStockRequest{
			ProductId: int32(item.ProductID),
			Quantity:  item.Quantity,
		})
		if err != nil || !resp.Available {
			return fmt.Errorf("product %d out of stock", item.ProductID)
		}
	}
	return uc.repo.Create(order)
}

func (uc *OrderUsecase) GetByID(id int) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *OrderUsecase) UpdateStatus(id int, status domain.OrderStatus) error {
	return uc.repo.UpdateStatus(id, status)
}

func (uc *OrderUsecase) ListByUser(userID int) ([]domain.Order, error) {
	return uc.repo.ListByUser(userID)
}

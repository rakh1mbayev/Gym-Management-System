package usecase

import "order_service/internal/domain"

type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(r domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: r}
}

func (uc *OrderUsecase) Create(order *domain.Order) error {
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

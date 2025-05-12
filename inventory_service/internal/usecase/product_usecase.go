package usecase

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/repository/postgres"
)

type ProductUsecase struct {
	repo postgres.ProductRepository
}

type ProductService interface {
	Create(ctx context.Context, p *domain.Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Product, error)
}

func NewProductUsecase(r postgres.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: r}
}

func (uc *ProductUsecase) Create(ctx context.Context, p *domain.Product) (int64, error) {
	return uc.repo.Create(ctx, p)
}

func (uc *ProductUsecase) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ProductUsecase) Update(ctx context.Context, p *domain.Product) error {
	return uc.repo.Update(ctx, p)
}

func (uc *ProductUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ProductUsecase) List(ctx context.Context) ([]domain.Product, error) {
	return uc.repo.List(ctx)
}

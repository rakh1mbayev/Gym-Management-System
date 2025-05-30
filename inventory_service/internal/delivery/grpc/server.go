package grpc

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/proto/inventorypb"
)

type InventoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	usecase usecase.ProductService
}

func NewInventoryServer(uc usecase.ProductService) *InventoryServer {
	return &InventoryServer{usecase: uc}
}

func (s *InventoryServer) CreateProduct(ctx context.Context, req *inventorypb.CreateProductRequest) (*inventorypb.Product, error) {
	product := &domain.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
	}

	id, err := s.usecase.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	product.ProductID = id

	return &inventorypb.Product{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}, nil
}

func (s *InventoryServer) GetProduct(ctx context.Context, req *inventorypb.GetProductRequest) (*inventorypb.Product, error) {
	product, err := s.usecase.GetByID(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &inventorypb.Product{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}, nil
}

func (s *InventoryServer) ListProducts(ctx context.Context, _ *inventorypb.ListProductsRequest) (*inventorypb.ListProductsResponse, error) {
	products, err := s.usecase.List(ctx)
	if err != nil {
		return nil, err
	}

	var resp inventorypb.ListProductsResponse
	for _, p := range products {
		resp.Products = append(resp.Products, &inventorypb.Product{
			ProductId:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
		})
	}

	return &resp, nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *inventorypb.UpdateProductRequest) (*inventorypb.Product, error) {
	product := &domain.Product{
		ProductID:   req.GetProductId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
	}

	if err := s.usecase.Update(ctx, product); err != nil {
		return nil, err
	}

	return &inventorypb.Product{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
	}, nil
}

func (s *InventoryServer) DeleteProduct(ctx context.Context, req *inventorypb.DeleteProductRequest) (*inventorypb.DeleteProductResponse, error) {
	if err := s.usecase.Delete(ctx, req.GetProductId()); err != nil {
		return nil, err
	}

	return &inventorypb.DeleteProductResponse{Success: true}, nil
}

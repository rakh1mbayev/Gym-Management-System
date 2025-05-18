package domain

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/proto/inventorypb"
	"google.golang.org/grpc"
)

type InventoryService interface {
	GetProduct(ctx context.Context, req *inventorypb.GetProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error)
	CreateProduct(ctx context.Context, req *inventorypb.CreateProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error)
	UpdateProduct(ctx context.Context, req *inventorypb.UpdateProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error)
	DeleteProduct(ctx context.Context, req *inventorypb.DeleteProductRequest, opts ...grpc.CallOption) (*inventorypb.DeleteProductResponse, error)
	ListProducts(ctx context.Context, req *inventorypb.ListProductsRequest, opts ...grpc.CallOption) (*inventorypb.ListProductsResponse, error)
}

type EventPublisher interface {
	PublishOrderCreated(data interface{}) error
}

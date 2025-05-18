package usecase

import (
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/proto/inventorypb"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

type mockOrderRepo struct {
	mock.Mock
}

func (m *mockOrderRepo) Create(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *mockOrderRepo) GetByID(id string) (*domain.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *mockOrderRepo) ListByUser(userID int64) ([]domain.Order, error) {
	args := m.Called(userID, int64(1))
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *mockOrderRepo) UpdateStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type mockNatsClient struct {
	mock.Mock
}

func (m *mockNatsClient) PublishOrderCreated(data interface{}) error {
	args := m.Called(data)
	return args.Error(0)
}

type mockInventoryClient struct {
	mock.Mock
}

func (m *mockInventoryClient) GetProduct(ctx context.Context, req *inventorypb.GetProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*inventorypb.Product), args.Error(1)
}

func (m *mockInventoryClient) CreateProduct(ctx context.Context, req *inventorypb.CreateProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*inventorypb.Product), args.Error(1)
}

func (m *mockInventoryClient) UpdateProduct(ctx context.Context, req *inventorypb.UpdateProductRequest, opts ...grpc.CallOption) (*inventorypb.Product, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*inventorypb.Product), args.Error(1)
}

func (m *mockInventoryClient) DeleteProduct(ctx context.Context, req *inventorypb.DeleteProductRequest, opts ...grpc.CallOption) (*inventorypb.DeleteProductResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*inventorypb.DeleteProductResponse), args.Error(1)
}

func (m *mockInventoryClient) ListProducts(ctx context.Context, req *inventorypb.ListProductsRequest, opts ...grpc.CallOption) (*inventorypb.ListProductsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*inventorypb.ListProductsResponse), args.Error(1)
}

func TestOrderUsecase_GetOrder(t *testing.T) {
	mockRepo := new(mockOrderRepo)
	mockNats := new(mockNatsClient)
	mockInventory := new(mockInventoryClient)

	uc := NewOrderUsecase(mockRepo, mockNats, mockInventory)

	// Arrange: expected order to be returned
	expectedOrder := &domain.Order{
		OrderID: "1",
		UserID:  101,
		Status:  "pending",
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 2},
		},
	}

	mockRepo.On("GetByID", "1").Return(expectedOrder, nil)

	// Act
	order, err := uc.GetOrder(context.Background(), "1")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order == nil || order.OrderID != "1" {
		t.Fatalf("expected order with ID 1, got %+v", order)
	}

	mockRepo.AssertExpectations(t)
	mockNats.AssertNotCalled(t, "PublishOrderCreated", mock.Anything)            // not relevant for GetOrder
	mockInventory.AssertNotCalled(t, "GetProduct", mock.Anything, mock.Anything) // not relevant for GetOrder
}

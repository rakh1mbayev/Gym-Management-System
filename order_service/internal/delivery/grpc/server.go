package grpc

import (
	"context"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/proto/orderpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

type OrderServiceServer struct {
	orderpb.UnimplementedOrderServiceServer
	UC usecase.OrderService
}

func NewOrderServiceServer(uc usecase.OrderService) *OrderServiceServer {
	return &OrderServiceServer{UC: uc}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	var items []domain.OrderItem
	var totalPrice float64

	for _, i := range req.GetItems() {
		item := domain.OrderItem{
			ProductID:    i.GetProductID(),
			Quantity:     int(i.GetQuantity()),
			PricePerItem: i.GetPricePerItem(),
		}
		items = append(items, item)
		totalPrice += float64(item.Quantity) * item.PricePerItem
	}

	orderID, err := s.UC.CreateOrder(ctx, req.GetUserId(), items, totalPrice)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		OrderId:    orderID,
		Status:     "pending", // Assuming new orders are pending by default
		TotalPrice: totalPrice,
	}, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderDetails, error) {
	order, err := s.UC.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	var items []*orderpb.OrderItem
	for _, i := range order.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId:    strconv.FormatInt(i.ProductID, 10), // if needed
			Quantity:     int32(i.Quantity),
			PricePerItem: float32(i.PricePerItem),
		})
	}

	return &orderpb.OrderDetails{
		OrderId:    order.OrderID,
		UserId:     order.UserID,
		Items:      items,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  timestamppb.New(order.CreatedAt),
		UpdatedAt:  timestamppb.New(order.UpdatedAt),
	}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *orderpb.OrderListRequest) (*orderpb.OrderListResponse, error) {
	orders, err := s.UC.ListOrders(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	var resp []*orderpb.OrderDetails
	for _, order := range orders {
		var items []*orderpb.OrderItem
		for _, i := range order.Items {
			items = append(items, &orderpb.OrderItem{
				ProductId:    strconv.Itoa(i.ProductID),
				Quantity:     int32(i.Quantity),
				PricePerItem: float32(i.PricePerItem),
			})
		}
		resp = append(resp, &orderpb.OrderDetails{
			OrderId:    order.OrderID,
			UserId:     order.UserID,
			Items:      items,
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
		})
	}

	return &orderpb.OrderListResponse{Orders: resp}, nil
}

func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	// Call the usecase to update the order status
	err := s.UC.UpdateOrderStatus(ctx, req.GetOrderId(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	// You can fetch updated order to return status & price or just return what you received
	return &orderpb.OrderResponse{
		OrderId: req.GetOrderId(),
		Status:  req.GetStatus(),
	}, nil
}

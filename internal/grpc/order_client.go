package grpc

import (
	"Gym-Management-System/order_service/proto/orderpb"
	"google.golang.org/grpc"
)

type OrderGRPCClient interface {
	orderpb.OrderServiceClient
}

func NewOrderGRPCClient(conn *grpc.ClientConn) OrderGRPCClient {
	return orderpb.NewOrderServiceClient(conn)
}

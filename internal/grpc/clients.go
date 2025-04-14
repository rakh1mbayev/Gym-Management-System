package grpc_clients

import (
	"log"

	"github.com/rakh1mbayev/Gym-Management-System/proto/userpb"
	"google.golang.org/grpc"
)

var UserClient userpb.UserServiceClient

func InitClients() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user_service: %v", err)
	}
	UserClient = userpb.NewUserServiceClient(conn)
}

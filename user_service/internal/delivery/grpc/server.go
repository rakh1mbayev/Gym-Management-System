package grpc

import (
	"Gym-Management-System/proto/userpb"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/internal/usecase"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	UC usecase.Usecase
}

func (s *Server) RegisterUser(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	id, err := s.UC.RegisterUser(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &userpb.UserResponse{
		Id:      id,
		Message: "User registered successfully",
	}, nil
}

// Implement AuthenticateUser and GetUserProfile similarly

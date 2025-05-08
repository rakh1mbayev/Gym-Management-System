package grpc

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/proto/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	Usecase usecase.UserService
}

const jwtSecret = "superSecret"

func NewUserServiceServer(uc usecase.UserService) *UserServiceServer {
	return &UserServiceServer{Usecase: uc}
}

func (s *UserServiceServer) RegisterUser(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	user := &domain.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	}
	err := s.Usecase.Register(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "registration failed: %v", err)
	}

	return &userpb.CreateResponse{UserId: user.ID}, nil
}

func (s *UserServiceServer) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	user, err := s.Usecase.Authenticate(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	// Build JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign it
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not sign token: %v", err)
	}

	return &userpb.AuthResponse{Token: signed}, nil
}

func (s *UserServiceServer) GetUserProfile(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	user, err := s.Usecase.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &userpb.GetResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *UserServiceServer) ConfirmEmail(ctx context.Context, req *userpb.ConfirmEmailRequest) (*userpb.ConfirmEmailResponse, error) {
	err := s.Usecase.ConfirmEmail(ctx, req.Token)
	if err != nil {
		return &userpb.ConfirmEmailResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &userpb.ConfirmEmailResponse{
		Success: true,
		Message: "Email confirmed successfully!",
	}, nil
}

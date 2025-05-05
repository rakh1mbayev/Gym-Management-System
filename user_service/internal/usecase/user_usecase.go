package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/delivery/grpc_client"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo       domain.UserRepository
	mailClient *grpc_client.MailClient
}

func NewUserUsecase(repo domain.UserRepository, mailClient *grpc_client.MailClient) *UserUsecase {
	return &UserUsecase{
		repo:       repo,
		mailClient: mailClient,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, user *domain.User) error {
	err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// Send confirmation email
	subject := "Welcome to Gym Management!"
	body := fmt.Sprintf("Hello %s! Please confirm your email address.", user.Name)

	return uc.mailClient.SendConfirmationEmail(user.Email, subject, body)
}

func (uc *UserUsecase) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (uc *UserUsecase) GetProfile(ctx context.Context, id int64) (*domain.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

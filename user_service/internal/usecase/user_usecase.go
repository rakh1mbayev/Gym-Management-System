package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/delivery/grpc_client"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type UserUsecase struct {
	repo       domain.UserRepository
	mailClient *grpc_client.MailClient
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	GetProfile(ctx context.Context, id int64) (*domain.User, error)
	ConfirmEmail(ctx context.Context, token string) error
}

func NewUserUsecase(repo domain.UserRepository, mailClient *grpc_client.MailClient) *UserUsecase {
	return &UserUsecase{
		repo:       repo,
		mailClient: mailClient,
	}
}

// Example function to generate a random token
func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)
	for i := range token {
		token[i] = letters[rand.Intn(len(letters))]
	}
	return string(token)
}

func (uc *UserUsecase) Register(ctx context.Context, user *domain.User) error {
	token := generateToken()
	user.ConfirmationToken = token
	user.IsConfirmed = false

	err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	err = uc.repo.SetConfirmationToken(ctx, user.ID, token)
	if err != nil {
		return err
	}

	confirmationURL := fmt.Sprintf("http://localhost:8083/confirm-email?token=%s", token)
	body := fmt.Sprintf("Click the link to confirm your email:\n%s", confirmationURL)

	return uc.mailClient.SendConfirmationEmail(user.Email, "Confirm your email", body)

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

func (uc *UserUsecase) ConfirmEmail(ctx context.Context, token string) error {
	user, err := uc.repo.GetUserByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	return uc.repo.MarkEmailConfirmed(ctx, user.ID)
}

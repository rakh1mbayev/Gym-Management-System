package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
)

type UserUsecase struct {
	repo       domain.UserRepository
	mailClient domain.MailService
	nats       domain.EventPublisher
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	GetProfile(ctx context.Context, id int64) (*domain.User, error)
	ConfirmEmail(ctx context.Context, token string) error
}

func NewUserUsecase(repo domain.UserRepository, mailClient domain.MailService, nats domain.EventPublisher) *UserUsecase {
	return &UserUsecase{
		repo:       repo,
		mailClient: mailClient,
		nats:       nats,
	}
}

// Example function to generate a random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (uc *UserUsecase) Register(ctx context.Context, user *domain.User) error {
	// Generate secure token
	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	// Save user first
	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	// Save token to DB
	if err := uc.repo.SetConfirmationToken(ctx, user.UserID, token); err != nil {
		return fmt.Errorf("failed to store confirmation token: %w", err)
	}

	if err := uc.nats.PublishUserRegistered(user.Email); err != nil {
		return fmt.Errorf("failed to publish user registered event: %w", err)
	}

	// Prepare email
	subject := "Please confirm your email"
	confirmationLink := fmt.Sprintf("http://localhost:8080/confirm?token=%s", token)
	body := fmt.Sprintf("Hello %s!\n\nPlease confirm your email by clicking the following link:\n%s", user.Name, confirmationLink)

	log.Println("NATS Publish: sending to info from user to email to confirm registration")
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

	if !user.IsConfirmed {
		return nil, errors.New("user is not confirmed")
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
	return uc.repo.MarkEmailConfirmed(ctx, user.UserID)
}

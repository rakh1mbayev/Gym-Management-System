package usercase

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"user_service/internal"
)

type Repo interface {
	CreateUser(ctx context.Context, username, email, hashedPassword string) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*internal.User, error)
	GetUserByID(ctx context.Context, id string) (*internal.User, error)
}

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) Usecase {
	return Usecase{repo: repo}
}

func (uc Usecase) RegisterUser(ctx context.Context, username, password, email string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return uc.repo.CreateUser(ctx, username, email, string(hashedPassword))
}

func (uc Usecase) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	user, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (uc Usecase) GetUserProfile(ctx context.Context, id string) (*internal.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

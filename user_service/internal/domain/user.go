package domain

import "golang.org/x/net/context"

type User struct {
	UserID            int64  `json:"user_id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	ConfirmationToken string `json:"confirmation_token"`
	IsConfirmed       bool   `json:"is_confirmed"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	SetConfirmationToken(ctx context.Context, userID int64, token string) error
	ConfirmUserByToken(ctx context.Context, token string) error
	MarkEmailConfirmed(ctx context.Context, id int64) error
	GetUserByToken(ctx context.Context, token string) (*User, error)
}

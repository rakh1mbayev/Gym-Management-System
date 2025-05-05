package domain

import "golang.org/x/net/context"

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Phone    string
	Role     string
}

type UserService interface {
	Register(ctx context.Context, user *User) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
	GetProfile(ctx context.Context, id int64) (*User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
}

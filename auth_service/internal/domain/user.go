package domain

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string // stored as a hashed password
	Phone    string
	Role     Role
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}

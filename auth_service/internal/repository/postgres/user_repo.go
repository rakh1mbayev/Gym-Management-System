package postgres

import (
	"auth_service/internal/domain"
	"database/sql"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *domain.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (name, email, password, phone, role) VALUES ($1, $2, $3, $4, $5)",
		user.Name, user.Email, user.Password, user.Phone, user.Role,
	)
	return err
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(
		"SELECT user_id, name, email, password, phone, role FROM users WHERE email = $1",
		email,
	)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"user_service/internal"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, username, email, hashedPassword string) (string, error) {
	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`, id, username, email, hashedPassword)

	return id, err
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*internal.User, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password FROM users WHERE username=$1
	`, username)

	var u internal.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*internal.User, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password FROM users WHERE id=$1
	`, id)

	var u internal.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

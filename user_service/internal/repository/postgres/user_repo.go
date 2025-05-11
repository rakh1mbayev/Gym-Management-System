package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING user_id`
	err := r.DB.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.Role).Scan(&user.UserID)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT user_id, name, email, password, role, confirmation_token, is_confirmed FROM users WHERE email=$1`
	row := r.DB.QueryRowContext(ctx, query, email)

	var user domain.User
	var confirmationToken sql.NullString // Handle nullable token

	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role, &confirmationToken, &user.IsConfirmed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	// If the confirmation token is not NULL, assign it
	if confirmationToken.Valid {
		user.ConfirmationToken = confirmationToken.String
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT user_id, name, email, password, role FROM users WHERE user_id=$1`
	row := r.DB.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	query := `UPDATE users SET confirmation_token = $1 WHERE user_id = $2`
	_, err := r.DB.ExecContext(ctx, query, token, userID)
	return err
}

func (r *UserRepository) ConfirmUserByToken(ctx context.Context, token string) error {
	query := `UPDATE users SET is_confirmed = true WHERE confirmation_token = $1`
	res, err := r.DB.ExecContext(ctx, query, token)
	fmt.Println(token, err)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("invalid or expired token")
	}
	return nil
}

func (r *UserRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	query := `SELECT user_id, name, email, password, role, is_confirmed FROM users WHERE confirmation_token = $1`

	row := r.DB.QueryRowContext(ctx, query, token)

	var user domain.User
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsConfirmed)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) MarkEmailConfirmed(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_confirmed = TRUE, confirmation_token = NULL WHERE user_id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

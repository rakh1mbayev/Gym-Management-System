package repository_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	dns := "host=localhost port=5432 user=postgres password=12345678 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dns)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	_, _ = db.Exec("DROP TABLE IF EXISTS users")
	_, err = db.Exec(`CREATE TABLE users (
		user_id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	return db
}

func TestGetUserByID_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	// Insert a test user
	insertQuery := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING user_id`
	var userID int64
	err := db.QueryRow(insertQuery, "Test User", "test@gmail.com", "password", "member").Scan(&userID)
	assert.NoError(t, err)

	// Fetch the user by ID
	user, err := repo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "member", user.Role)
}

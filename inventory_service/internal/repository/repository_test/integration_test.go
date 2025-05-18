package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func setupDB(t *testing.T) *sql.DB {
	dns := "host=localhost port=5432 user=postgres password=12345678 dbname=test sslmode=disable"
	db, err := sql.Open("postgres", dns)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	_, _ = db.Exec("DROP TABLE IF EXISTS products")
	_, err = db.Exec(`CREATE TABLE products (
		product_id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		product_description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		stock INT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	return db
}

func TestProductRepo_List(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	repo := postgres.NewProductRepository(db, nil)

	insertQuery := `INSERT INTO products (name, product_description, price, stock) VALUES ($1, $2, $3, $4) RETURNING product_id`
	var id int64
	err := db.QueryRow(insertQuery, "Dumbbell", "20kg weight", 49.99, 10).Scan(&id)
	assert.NoError(t, err)
	err = db.QueryRow(insertQuery, "Bar", "20kg weight", 10.99, 5).Scan(&id)
	assert.NoError(t, err)

	products, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Dumbbell", products[0].Name)
	assert.Equal(t, "Bar", products[1].Name)
}

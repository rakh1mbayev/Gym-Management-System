package repository_test

import (
	"context"
	_ "database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/domain"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/repository/postgres"
)

func TestProductRepo_Create(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewProductRepository(db, nil)

	product := &domain.Product{
		Name:        "Dumbbell",
		Description: "20kg weight",
		Price:       49.99,
		Stock:       10,
	}
	expectedID := int64(1)

	mock.ExpectQuery("INSERT INTO products").
		WithArgs(product.Name, product.Description, product.Price, product.Stock).
		WillReturnRows(sqlmock.NewRows([]string{"product_id"}).AddRow(expectedID))

	id, err := repo.Create(context.Background(), product)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

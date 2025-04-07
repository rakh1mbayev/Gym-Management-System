package postgres

import (
	"context"
	"database/sql"
	"inventory_service/internal/domain"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO products (name, product_description, price) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Price)
	return err
}

func (r *productRepo) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT product_id, name, product_description, price FROM products WHERE product_id = $1", id)

	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	return &p, err
}

func (r *productRepo) Update(ctx context.Context, p *domain.Product) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE products SET name=$1, product_description=$2, price=$3 WHERE product_id=$4",
		p.Name, p.Description, p.Price, p.ID)
	return err
}

func (r *productRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM products WHERE product_id = $1", id)
	return err
}

func (r *productRepo) List(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT product_id, name, product_description, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

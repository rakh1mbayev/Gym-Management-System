package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/domain"
	"github.com/redis/go-redis/v9"
	"log"
)

type productRepo struct {
	db    *sql.DB
	cache *redis.Client
}

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Product, error)
	DecreaseStock(ctx context.Context, productID int64, quantity int) error
}

func NewProductRepository(db *sql.DB, cache *redis.Client) ProductRepository {
	return &productRepo{db: db, cache: cache}
}

func (r *productRepo) Update(ctx context.Context, p *domain.Product) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE products SET name=$1, product_description=$2, price=$3, stock=$4 WHERE product_id=$5",
		p.Name, p.Description, p.Price, p.Stock, p.ProductID)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", p.ProductID)
	_ = r.cache.Del(ctx, cacheKey).Err()

	return nil
}

func (r *productRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", id)
	_ = r.cache.Del(ctx, cacheKey).Err()

	return nil
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO products (name, product_description, price, stock) VALUES ($1, $2, $3, $4) RETURNING product_id",
		p.Name, p.Description, p.Price, p.Stock).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *productRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	val, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedProduct domain.Product
		if err := json.Unmarshal([]byte(val), &cachedProduct); err == nil {
			log.Println("CACHE HIT")
			return &cachedProduct, nil
		}
	}

	row := r.db.QueryRowContext(ctx,
		"SELECT product_id, name, product_description, price, stock FROM products WHERE product_id = $1", id)

	var p domain.Product
	err = row.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(p)
	_ = r.cache.Set(ctx, cacheKey, data, 0).Err()

	log.Println("CACHE MISS")
	return &p, nil
}

func (r *productRepo) List(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT product_id, name, product_description, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepo) DecreaseStock(ctx context.Context, productID int64, quantity int) error {
	// Make sure stock doesn't go negative
	_, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET stock = stock - $1
		WHERE product_id = $2 AND stock >= $1
	`, quantity, productID)

	if err != nil {
		return err
	}

	// Check if the update actually happened
	var updated bool
	err = r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM products WHERE product_id = $1 AND stock >= 0
		)
	`, productID).Scan(&updated)
	if err != nil {
		return err
	}
	if !updated {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}
	return nil
}

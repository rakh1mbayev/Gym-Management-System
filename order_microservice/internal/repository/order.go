package repository

import (
	"context"
	"database/sql"
)

type productRepo struct {
	db *sql.DB
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO products ...")
	return err
}

package postgres

import (
	"database/sql"
	"order_serivce/internal/domain"
	"time"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (user_id, name, price, status) VALUES ($1, $2, $3, $4) RETURNING order_id`
	err = tx.QueryRow(query, order.UserID, order.Name, order.Price, order.Status).Scan(&order.ID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`INSERT INTO order_items (order_id, product_id, quantity, price_per_item)
			VALUES ($1, $2, $3, $4)`, order.ID, item.ProductID, item.Quantity, item.PricePerItem)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepo) GetByID(id int) (*domain.Order, error) {
	row := r.db.QueryRow(`SELECT order_id, user_id, name, price, status, created_at, updated_at FROM orders WHERE order_id = $1`, id)

	var o domain.Order
	err := row.Scan(&o.ID, &o.UserID, &o.Name, &o.Price, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	itemsRows, err := r.db.Query(`SELECT order_item_id, product_id, quantity, price_per_item FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item domain.OrderItem
		item.OrderID = id
		err := itemsRows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.PricePerItem)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *orderRepo) UpdateStatus(id int, status domain.OrderStatus) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3`, status, time.Now(), id)
	return err
}

func (r *orderRepo) ListByUser(userID int) ([]domain.Order, error) {
	rows, err := r.db.Query(`SELECT order_id, name, price, status, created_at, updated_at FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		o.UserID = userID
		err := rows.Scan(&o.ID, &o.Name, &o.Price, &o.Status, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

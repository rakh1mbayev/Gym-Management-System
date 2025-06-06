package postgres

import (
	"database/sql"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/domain"
	"log"
	"time"
)

type orderRepo struct {
	db *sql.DB
}

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status string) error
	ListByUser(userID int64) ([]domain.Order, error)
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order *domain.Order) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Failed to insert order: %v\n", err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (user_id, total_price, status) 
	          VALUES ($1, $2, $3) RETURNING order_id`
	err = tx.QueryRow(query, order.UserID, order.TotalPrice, order.Status).Scan(&order.OrderID)
	if err != nil {
		log.Printf("Failed to insert order: %v\n", err)
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`INSERT INTO order_items (order_id, product_id, quantity, price_per_item)
			VALUES ($1, $2, $3, $4)`, order.OrderID, item.ProductID, item.Quantity, item.PricePerItem)
		if err != nil {
			log.Printf("Failed to insert order: %v\n", err)
			return err
		}
	}

	log.Printf("Successfully created order %s with %d items\n", order.OrderID, len(order.Items))
	return tx.Commit()
}

func (r *orderRepo) GetByID(id string) (*domain.Order, error) {
	row := r.db.QueryRow(`SELECT order_id, user_id, total_price, status 
	                       FROM orders WHERE order_id = $1`, id)

	var order domain.Order
	err := row.Scan(&order.OrderID, &order.UserID, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}

	// Fetch order items
	itemsRows, err := r.db.Query(`SELECT product_id, quantity, price_per_item 
	                               FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item domain.OrderItem
		err := itemsRows.Scan(&item.ProductID, &item.Quantity, &item.PricePerItem)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *orderRepo) UpdateStatus(id string, status string) error {
	// Update order status
	_, err := r.db.Exec(`UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3`,
		status, time.Now(), id)
	return err
}

func (r *orderRepo) ListByUser(userID int64) ([]domain.Order, error) {
	rows, err := r.db.Query(`SELECT order_id, total_price, status 
	                         FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.OrderID, &order.TotalPrice, &order.Status)
		if err != nil {
			return nil, err
		}

		// Fetch the items for each order
		itemsRows, err := r.db.Query(`SELECT product_id, quantity, price_per_item 
		                               FROM order_items WHERE order_id = $1`, order.OrderID)
		if err != nil {
			return nil, err
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			var item domain.OrderItem
			err := itemsRows.Scan(&item.ProductID, &item.Quantity, &item.PricePerItem)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

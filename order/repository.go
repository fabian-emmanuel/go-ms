package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database::{%s}::%w", url, err)
	}

	// Verify database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database::{%s}::%w", url, err)
	}

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	if r.db == nil {
		log.Println("Warning: Attempted to close a nil database connection")
		return
	}

	if err := r.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) CreateOrder(ctx context.Context, order Order) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback only if commit hasn't occurred
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // Rollback on panic
			panic(p)          // Re-throw panic
		} else if err != nil {
			_ = tx.Rollback() // Rollback if there was an error
		}
	}()

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_amount) VALUES($1, $2, $3, $4)",
		order.ID, order.CreatedAt, order.AccountId, order.TotalAmount,
	)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	for _, orderedProduct := range order.Products {
		_, err = stmt.ExecContext(ctx, order.ID, orderedProduct.ID, orderedProduct.Quantity)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return err
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id, o.created_at, o.account_id, o.total_amount::money::numeric::float8, 
		        op.product_id, op.quantity 
		 FROM orders o 
		 JOIN ordered_products op ON o.id = op.order_id 
		 WHERE o.account_id = $1 
		 ORDER BY o.id`,
		accountId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	// Grouping orders
	ordersMap := make(map[string]*Order)
	var orders []*Order

	for rows.Next() {
		var orderID, accountID string
		var createdAt time.Time
		var totalAmount float64
		var productID string
		var quantity uint32

		if err := rows.Scan(&orderID, &createdAt, &accountID, &totalAmount, &productID, &quantity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Check if order already exists in the map
		if _, exists := ordersMap[orderID]; !exists {
			ordersMap[orderID] = &Order{
				ID:          orderID,
				CreatedAt:   createdAt,
				AccountId:   accountID,
				TotalAmount: totalAmount,
				Products:    []OrderedProduct{},
			}
			orders = append(orders, ordersMap[orderID])
		}

		// Append the product to the order
		ordersMap[orderID].Products = append(ordersMap[orderID].Products, OrderedProduct{
			ID:       productID,
			Quantity: quantity,
		})
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return orders, nil
}

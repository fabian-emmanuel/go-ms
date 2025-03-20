package account

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip, take uint64) ([]*Account, error)
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

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	a := &Account{}

	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip, take uint64) ([]*Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var accounts []*Account

	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, a)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

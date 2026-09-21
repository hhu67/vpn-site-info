package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}

func InitSchema(db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS auth (
		id SERIAL PRIMARY KEY,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vpn_links (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		link VARCHAR(1024) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(context.Background(), query)
	return err
}

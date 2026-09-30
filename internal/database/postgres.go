package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func CreateTables(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role VARCHAR(30) NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS service_requests (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			phone VARCHAR(30) NOT NULL,
			address VARCHAR(255) NOT NULL,
			problem TEXT NOT NULL,
			status VARCHAR(30) NOT NULL DEFAULT 'new',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}

	return nil
}

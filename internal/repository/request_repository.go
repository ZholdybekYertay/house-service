package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"service/internal/model"
)

var ErrNotFound = errors.New("request not found")

type RequestRepository struct {
	db *pgxpool.Pool
}

func NewRequestRepository(db *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(ctx context.Context, input model.CreateRequestInput) (model.ServiceRequest, error) {
	const query = `
		INSERT INTO service_requests (name, phone, address, problem)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, phone, address, problem, status, created_at
	`

	var item model.ServiceRequest

	err := r.db.QueryRow(
		ctx,
		query,
		input.Name,
		input.Phone,
		input.Address,
		input.Problem,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Address,
		&item.Problem,
		&item.Status,
		&item.CreatedAt,
	)
	if err != nil {
		return model.ServiceRequest{}, fmt.Errorf("create request: %w", err)
	}

	return item, nil
}

func (r *RequestRepository) List(ctx context.Context) ([]model.ServiceRequest, error) {
	const query = `
		SELECT id, name, phone, address, problem, status, created_at
		FROM service_requests
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list requests: %w", err)
	}
	defer rows.Close()

	items := make([]model.ServiceRequest, 0)

	for rows.Next() {
		var item model.ServiceRequest

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Address,
			&item.Problem,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan request: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate requests: %w", err)
	}

	return items, nil
}

func (r *RequestRepository) UpdateStatus(ctx context.Context, id int64, status string) (model.ServiceRequest, error) {
	const query = `
		UPDATE service_requests
		SET status = $2
		WHERE id = $1
		RETURNING id, name, phone, address, problem, status, created_at
	`

	var item model.ServiceRequest

	err := r.db.QueryRow(ctx, query, id, status).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Address,
		&item.Problem,
		&item.Status,
		&item.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceRequest{}, ErrNotFound
	}
	if err != nil {
		return model.ServiceRequest{}, fmt.Errorf("update request status: %w", err)
	}

	return item, nil
}

func (r *RequestRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM service_requests WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

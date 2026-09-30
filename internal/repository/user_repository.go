package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"service/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (model.User, error) {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, role
	`

	var user model.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const query = `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

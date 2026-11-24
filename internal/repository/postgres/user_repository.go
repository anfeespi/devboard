package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfeespi/devboard/internal/domain"
	postgres "github.com/anfeespi/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements use case using postgresql
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *postgres.Queries
}

// NewUserRepository constructor
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: postgres.New(pool),
	}
}

// Create create user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.queries.CreateUser(ctx, postgres.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	})

	if err != nil {
		return fmt.Errorf("inserting user in bd: %w", err)
	}

	user.ID = result.ID.String()
	user.CreatedAt = result.CreatedAt.Time
	user.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// GetByID search a user by its id
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetUserByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consulting user %s: %w", id, err)
	}

	return toDomainUser(result), nil
}

// GetByEmail obtain a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("email %s: %w", email, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consulting email %s: %w", email, err)
	}

	return toDomainUser(result), nil
}

// toDomainUser maps the model sqlc to the domain user
func toDomainUser(u postgres.User) *domain.User {
	return &domain.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

// Package memory store things in memory
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/anfeespi/devboard/internal/domain"
)

// UserRepository is an adapter that implement userrepository
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewUserRepository create a new user repository in memory
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// Create create a new user
func (r *UserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("user_%d", len(r.users)+1)
	r.users[user.ID] = user
	return nil
}

// GetByID obtain a user by its ID
func (r *UserRepository) GetByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]

	if !ok {
		return nil, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
	}

	return user, nil
}

// GetByEmail obtain a user by its email
func (r *UserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, fmt.Errorf("user %s: %w", email, domain.ErrNotFound)
}

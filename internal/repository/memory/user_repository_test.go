package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/anfeespi/devboard/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Andres",
		Email: "andres@devboard.dev",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() returned an unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Fatalf("Create() doesn't assign ID to the user")
	}

	if user.ID != "user_1" {
		t.Fatalf("Create() assigned %q, expected %q", user.ID, "user_1")
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "user_999")
	if err == nil {
		t.Fatalf("GetByID() must returned an error")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByID() returned %v, expected ErrNotFound", err)
	}
}

func TestUserRepository_GetByID_Successful(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Andres",
		Email: "andres@devboard.dev",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("cannot prepare the user: %v", err)
	}

	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() returned an unexpected error: %v", err)
	}

	if found.ID != user.ID {
		t.Fatalf("GetByID() returned ID %q, expected %q", found.ID, user.ID)
	}

	if found.Email != user.Email {
		t.Errorf("GetByID() returned email %q, expected %q", found.Email, user.Email)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Andres",
		Email: "andres@devboard.com",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("cannot prepare the user: %v", err)
	}

	found, err := repo.GetByEmail(ctx, "andres@devboard.com")
	if err != nil {
		t.Fatalf("GetByEmail() returned an unexpected error: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("GetByID() returned email %q, expected %q", found.Email, user.Email)
	}
}

package domain

import (
	"time"
)

// User represent the user of the system
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new user
func NewUser(email, name, passwordHash string) *User {
	now := time.Now()

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsValid check if the user is valid
func (u *User) IsValid() bool {
	return u.Email != "" && u.Name != ""
}

package domain

import "context"

// UserRepository defines the contract for user database operations.
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

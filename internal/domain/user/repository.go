package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, us *User) (*User, error)
	Update(ctx context.Context, us *User) (*User, error)
	Delete(ctx context.Context, userID int64) error
	GetByID(ctx context.Context, userID int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

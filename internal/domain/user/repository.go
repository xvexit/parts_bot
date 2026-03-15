package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, us *User) error
	Delete(ctx context.Context, userID int64) error
	GetByID(ctx context.Context, userID int64) (*User, error)
	GetByTgID(ctx context.Context, userTgID int64) (*User, error)
}
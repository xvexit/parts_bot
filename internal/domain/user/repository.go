package user

import "context"

type Repository interface{
	Save(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int) error
	GetByID(ctx context.Context, userID int) (User, error)
	GetByTgID(ctx context.Context, userTgID int) (User, error)
}
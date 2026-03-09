package cart

import "context"

type Repository interface {
	Save(ctx context.Context, cart *Cart) error
	GetByUserID(ctx context.Context, userID int64) (*Cart, error)
	Delete(ctx context.Context, userID int64) error
}
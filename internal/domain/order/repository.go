package order

import "context"

type Repository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	ListByUserID(ctx context.Context, userID int64) ([]*Order, error)
}

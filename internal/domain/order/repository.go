package order

import "context"

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	ListByUserID(ctx context.Context, userID int64) ([]Order, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus) error
}

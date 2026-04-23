package payment

import "context"

type Repository interface {
	Create(ctx context.Context, payment *Payment) (*Payment, error)
	Update(ctx context.Context, payment *Payment) (*Payment, error)
	GetByID(ctx context.Context, id int64) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID int64) (*Payment, error)
	GetByProviderTxnID(ctx context.Context, providerTxnID string) (*Payment, error)
	ListByOrderID(ctx context.Context, orderID int64) ([]Payment, error)
}
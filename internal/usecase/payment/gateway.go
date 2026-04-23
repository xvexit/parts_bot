package payment

import "context"

type Gateway interface {
	CreatePayment(ctx context.Context, input CreatePaymentGatewayInput) (*CreatePaymentGatewayResult, error)
	GetPaymentStatus(ctx context.Context, providerTxnID string) (string, error)
}
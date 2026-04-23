package payment

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MockGateway struct {
	baseURL string
}

func NewMockGateway(baseURL string) *MockGateway {
	return &MockGateway{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
	}
}

func (g *MockGateway) CreatePayment(ctx context.Context, input CreatePaymentGatewayInput) (*CreatePaymentGatewayResult, error) {
	_ = ctx

	txnID := fmt.Sprintf("mock_%d_%d", input.OrderID, time.Now().UnixNano())
	url := fmt.Sprintf(
		"%s/payment-result.html?order_id=%d&amount=%d&provider_txn_id=%s",
		g.baseURL,
		input.OrderID,
		input.Amount,
		txnID,
	)

	return &CreatePaymentGatewayResult{
		ProviderTxnID: txnID,
		PaymentURL:    url,
		Status:        "pending",
	}, nil
}

func (g *MockGateway) GetPaymentStatus(ctx context.Context, providerTxnID string) (string, error) {
	_ = ctx
	_ = providerTxnID
	return "pending", nil
}

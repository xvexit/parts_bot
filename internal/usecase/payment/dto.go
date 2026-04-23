package payment

type CreatePaymentInput struct {
	UserID    int64
	OrderID   int64
	ReturnURL string
}

type CreatePaymentGatewayInput struct {
	OrderID     int64
	Amount      int64
	Description string
	ReturnURL   string
}

type CreatePaymentGatewayResult struct {
	ProviderTxnID string
	PaymentURL    string
	Status        string
}
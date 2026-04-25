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

type YooKassaWebhookPayload struct {
	Type   string `json:"type"`
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Paid   bool   `json:"paid"`
	} `json:"object"`
}

type yooKassaCreatePaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool   `json:"capture"`
	Description  string `json:"description,omitempty"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type yooKassaPaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	Paid         bool   `json:"paid"`
	Test         bool   `json:"test"`
	CreatedAt    string `json:"created_at"`
	Cancellation *struct {
		Party  string `json:"party"`
		Reason string `json:"reason"`
	} `json:"cancellation_details,omitempty"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation *struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type yooKassaErrorResponse struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Parameter   string `json:"parameter,omitempty"`
}
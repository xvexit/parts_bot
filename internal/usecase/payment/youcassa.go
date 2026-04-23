package payment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type YooKassaGateway struct {
	shopID     string
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

func NewYooKassaGateway(shopID, secretKey string) *YooKassaGateway {
	return &YooKassaGateway{
		shopID:    strings.TrimSpace(shopID),
		secretKey: strings.TrimSpace(secretKey),
		baseURL:   "https://api.yookassa.ru/v3",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
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

func (g *YooKassaGateway) CreatePayment(ctx context.Context, input CreatePaymentGatewayInput) (*CreatePaymentGatewayResult, error) {
	reqBody := yooKassaCreatePaymentRequest{}
	reqBody.Amount.Value = formatAmountRub(input.Amount)
	reqBody.Amount.Currency = "RUB"
	reqBody.Capture = true
	reqBody.Description = input.Description
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = strings.TrimSpace(input.ReturnURL)
	reqBody.Metadata = map[string]string{
		"order_id": fmt.Sprintf("%d", input.OrderID),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		g.baseURL+"/payments",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewString())
	req.Header.Set("Authorization", g.basicAuthHeader())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr yooKassaErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Description != "" {
			return nil, fmt.Errorf("yookassa create payment failed: %s", apiErr.Description)
		}
		return nil, fmt.Errorf("yookassa create payment failed: http %d", resp.StatusCode)
	}

	var paymentResp yooKassaPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, err
	}

	var paymentURL string
	if paymentResp.Confirmation != nil {
		paymentURL = paymentResp.Confirmation.ConfirmationURL
	}

	return &CreatePaymentGatewayResult{
		ProviderTxnID: paymentResp.ID,
		PaymentURL:    paymentURL,
		Status:        normalizeYooKassaStatus(paymentResp.Status),
	}, nil
}

func (g *YooKassaGateway) GetPaymentStatus(ctx context.Context, providerTxnID string) (string, error) {
	providerTxnID = strings.TrimSpace(providerTxnID)
	if providerTxnID == "" {
		return "", fmt.Errorf("provider transaction id is required")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		g.baseURL+"/payments/"+providerTxnID,
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", g.basicAuthHeader())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr yooKassaErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Description != "" {
			return "", fmt.Errorf("yookassa get payment failed: %s", apiErr.Description)
		}
		return "", fmt.Errorf("yookassa get payment failed: http %d", resp.StatusCode)
	}

	var paymentResp yooKassaPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return "", err
	}

	return normalizeYooKassaStatus(paymentResp.Status), nil
}

func (g *YooKassaGateway) basicAuthHeader() string {
	raw := g.shopID + ":" + g.secretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

func formatAmountRub(amount int64) string {
	rubles := amount / 100
	kopecks := amount % 100
	return fmt.Sprintf("%d.%02d", rubles, kopecks)
}

func normalizeYooKassaStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending":
		return "pending"
	case "waiting_for_capture":
		return "succeeded"
	case "succeeded":
		return "succeeded"
	case "canceled":
		return "canceled"
	default:
		return "pending"
	}
}
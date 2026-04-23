package payment

import (
	"strings"
	"time"

	"partsBot/pkg/errors"
)

type Payment struct {
	id            int64
	orderID       int64
	amount        int64
	provider      string
	providerTxnID *string
	paymentURL    *string
	status        Status
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPayment(orderID, amount int64, provider string) (*Payment, error) {
	if orderID <= 0 {
		return nil, errors.ErrOrderEmpty
	}

	if amount <= 0 {
		return nil, errors.ErrAmountCanNotBeNull
	}

	provider = strings.TrimSpace(provider)
	if provider == "" {
		return nil, errors.ErrInvalidPaymentProvider
	}

	now := time.Now()

	return &Payment{
		orderID:    orderID,
		amount:     amount,
		provider:   provider,
		status:     StatusPending,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func RestorePayment(
	id int64,
	orderID int64,
	amount int64,
	provider string,
	providerTxnID *string,
	paymentURL *string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) *Payment {
	var st Status

	switch status {
	case string(StatusPending):
		st = StatusPending
	case string(StatusSucceeded):
		st = StatusSucceeded
	case string(StatusFailed):
		st = StatusFailed
	case string(StatusCanceled):
		st = StatusCanceled
	default:
		st = StatusFailed
	}

	provider = strings.TrimSpace(provider)

	return &Payment{
		id:            id,
		orderID:       orderID,
		amount:        amount,
		provider:      provider,
		providerTxnID: normalizeOptionalString(providerTxnID),
		paymentURL:    normalizeOptionalString(paymentURL),
		status:        st,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (p *Payment) SetID(id int64) {
	p.id = id
}

func (p *Payment) SetPaymentURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return errors.ErrInvalidPaymentURL
	}

	p.paymentURL = &url
	p.touch()
	return nil
}

func (p *Payment) SetProviderTxnID(txnID string) error {
	txnID = strings.TrimSpace(txnID)
	if txnID == "" {
		return errors.ErrInvalidProviderTxnID
	}

	p.providerTxnID = &txnID
	p.touch()
	return nil
}

func (p *Payment) MarkSucceeded() {
	p.status = StatusSucceeded
	p.touch()
}

func (p *Payment) MarkFailed() {
	p.status = StatusFailed
	p.touch()
}

func (p *Payment) MarkCanceled() {
	p.status = StatusCanceled
	p.touch()
}

func (p *Payment) touch() {
	p.updatedAt = time.Now()
}

func normalizeOptionalString(s *string) *string {
	if s == nil {
		return nil
	}

	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}

	return &v
}

func (p *Payment) ID() int64 {
	return p.id
}

func (p *Payment) OrderID() int64 {
	return p.orderID
}

func (p *Payment) Amount() int64 {
	return p.amount
}

func (p *Payment) Provider() string {
	return p.provider
}

func (p *Payment) ProviderTxnID() *string {
	return p.providerTxnID
}

func (p *Payment) PaymentURL() *string {
	return p.paymentURL
}

func (p *Payment) Status() string {
	return string(p.status)
}

func (p *Payment) StatusValue() Status {
	return p.status
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}
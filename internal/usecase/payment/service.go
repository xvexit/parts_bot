package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	dOrder "partsBot/internal/domain/order"
	dPayment "partsBot/internal/domain/payment"
)

type Service struct {
	repo      dPayment.Repository
	orderRepo dOrder.Repository
	gateway   Gateway
}

func NewService(
	repo dPayment.Repository,
	orderRepo dOrder.Repository,
	gateway Gateway,
) *Service {
	return &Service{
		repo:      repo,
		orderRepo: orderRepo,
		gateway:   gateway,
	}
}

// CreatePaymentForOrder создает новую попытку оплаты для заказа.
// Даже если по заказу уже были failed/canceled попытки, новая попытка допустима.
func (s *Service) CreatePaymentForOrder(
	ctx context.Context,
	input CreatePaymentInput,
) (*dPayment.Payment, error) {
	if input.UserID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if input.OrderID <= 0 {
		return nil, errors.New("invalid order id")
	}

	ord, err := s.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}

	if ord.UserID() != input.UserID {
		return nil, errors.New("access denied to order")
	}

	// Если заказ уже оплачен/подтвержден/доставлен, новую оплату создавать не надо.
	switch ord.Status() {
	case string(dOrder.OrderStatusConfirmed), string(dOrder.OrderStatusDelivered):
		return nil, errors.New("order is already paid or completed")
	}

	pay, err := dPayment.NewPayment(
		ord.ID(),
		ord.Total().Amount(),
		"mock", // потом подставишь "yookassa" / "tbank"
	)
	if err != nil {
		return nil, err
	}

	pay, err = s.repo.Create(ctx, pay)
	if err != nil {
		return nil, err
	}

	gwResult, err := s.gateway.CreatePayment(ctx, CreatePaymentGatewayInput{
		OrderID:     ord.ID(),
		Amount:      ord.Total().Amount(),
		Description: buildPaymentDescription(ord),
		ReturnURL:   strings.TrimSpace(input.ReturnURL),
	})
	if err != nil {
		pay.MarkFailed()
		pay, _ = s.repo.Update(ctx, pay)
		return nil, err
	}

	if err := pay.SetProviderTxnID(gwResult.ProviderTxnID); err != nil {
		return nil, err
	}

	if err := pay.SetPaymentURL(gwResult.PaymentURL); err != nil {
		return nil, err
	}

	switch normalizeGatewayStatus(gwResult.Status) {
	case string(dPayment.StatusSucceeded):
		pay.MarkSucceeded()
		_ = s.orderRepo.UpdateStatus(ctx, ord.ID(), dOrder.OrderStatusConfirmed)
	case string(dPayment.StatusFailed):
		pay.MarkFailed()
	case string(dPayment.StatusCanceled):
		pay.MarkCanceled()
	default:
		// pending — ничего не делаем
	}

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return nil, err
	}

	return pay, nil
}

func (s *Service) GetByID(ctx context.Context, paymentID int64) (*dPayment.Payment, error) {
	if paymentID <= 0 {
		return nil, errors.New("invalid payment id")
	}
	return s.repo.GetByID(ctx, paymentID)
}

// GetLastByOrderID возвращает последнюю попытку оплаты.
func (s *Service) GetLastByOrderID(ctx context.Context, userID, orderID int64) (*dPayment.Payment, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if orderID <= 0 {
		return nil, errors.New("invalid order id")
	}

	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if ord.UserID() != userID {
		return nil, errors.New("access denied to order")
	}

	return s.repo.GetByOrderID(ctx, orderID)
}

func (s *Service) ListByOrderID(ctx context.Context, userID, orderID int64) ([]dPayment.Payment, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if orderID <= 0 {
		return nil, errors.New("invalid order id")
	}

	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if ord.UserID() != userID {
		return nil, errors.New("access denied to order")
	}

	return s.repo.ListByOrderID(ctx, orderID)
}

// SyncStatusByProviderTxnID полезен для webhook или периодической сверки статуса.
func (s *Service) SyncStatusByProviderTxnID(ctx context.Context, providerTxnID string) (*dPayment.Payment, error) {
	providerTxnID = strings.TrimSpace(providerTxnID)
	if providerTxnID == "" {
		return nil, errors.New("provider transaction id is required")
	}

	pay, err := s.repo.GetByProviderTxnID(ctx, providerTxnID)
	if err != nil {
		return nil, err
	}

	status, err := s.gateway.GetPaymentStatus(ctx, providerTxnID)
	if err != nil {
		return nil, err
	}

	switch normalizeGatewayStatus(status) {
	case string(dPayment.StatusSucceeded):
		pay.MarkSucceeded()
		_ = s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusConfirmed)
	case string(dPayment.StatusFailed):
		pay.MarkFailed()
	case string(dPayment.StatusCanceled):
		pay.MarkCanceled()
	default:
		// pending
	}

	return s.repo.Update(ctx, pay)
}

// MarkSucceededByProviderTxnID — удобный метод под входящий callback.
func (s *Service) MarkSucceededByProviderTxnID(ctx context.Context, providerTxnID string) (*dPayment.Payment, error) {
	providerTxnID = strings.TrimSpace(providerTxnID)
	if providerTxnID == "" {
		return nil, errors.New("provider transaction id is required")
	}

	pay, err := s.repo.GetByProviderTxnID(ctx, providerTxnID)
	if err != nil {
		return nil, err
	}

	pay.MarkSucceeded()

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusConfirmed); err != nil {
		return nil, err
	}

	return pay, nil
}

func buildPaymentDescription(ord *dOrder.Order) string {
	return fmt.Sprintf("Оплата заказа #%d", ord.ID())
}

func normalizeGatewayStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending", "waiting_for_capture", "waiting", "":
		return string(dPayment.StatusPending)
	case "succeeded", "success", "paid":
		return string(dPayment.StatusSucceeded)
	case "failed", "error":
		return string(dPayment.StatusFailed)
	case "canceled", "cancelled":
		return string(dPayment.StatusCanceled)
	default:
		return string(dPayment.StatusPending)
	}
}
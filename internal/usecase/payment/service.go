package payment

import (
	"context"
	"encoding/json"
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

	switch ord.Status() {
	case string(dOrder.OrderStatusPaid),
		string(dOrder.OrderStatusConfirmed),
		string(dOrder.OrderStatusDelivered):
		return nil, errors.New("order is already paid or completed")
	}

	pay, err := dPayment.NewPayment(
		ord.ID(),
		ord.Total().Amount(),
		"yookassa",
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
		Description: fmt.Sprintf("Оплата заказа #%d", ord.ID()),
		ReturnURL:   strings.TrimSpace(input.ReturnURL),
	})
	if err != nil {
		pay.MarkFailed()
		_, _ = s.repo.Update(ctx, pay)
		return nil, err
	}

	if err := pay.SetProviderTxnID(gwResult.ProviderTxnID); err != nil {
		return nil, err
	}

	if err := pay.SetPaymentURL(gwResult.PaymentURL); err != nil {
		return nil, err
	}

	s.applyPaymentStatus(pay, gwResult.Status)

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return nil, err
	}

	if pay.Status() == string(dPayment.StatusSucceeded) {
		_ = s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusPaid)
	}

	return pay, nil
}

func (s *Service) GetLastByOrderID(ctx context.Context, userID, orderID int64) (*dPayment.Payment, error) {
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
	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if ord.UserID() != userID {
		return nil, errors.New("access denied to order")
	}

	return s.repo.ListByOrderID(ctx, orderID)
}

func (s *Service) SyncOrderPayment(ctx context.Context, userID, orderID int64) (*dPayment.Payment, error) {
	ord, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if ord.UserID() != userID {
		return nil, errors.New("access denied to order")
	}

	pay, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if pay.ProviderTxnID() == nil || strings.TrimSpace(*pay.ProviderTxnID()) == "" {
		return pay, nil
	}

	status, err := s.gateway.GetPaymentStatus(ctx, *pay.ProviderTxnID())
	if err != nil {
		return nil, err
	}

	s.applyPaymentStatus(pay, status)

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return nil, err
	}

	if pay.Status() == string(dPayment.StatusSucceeded) {
		_ = s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusPaid)
	}

	return pay, nil
}

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

	s.applyPaymentStatus(pay, status)

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return nil, err
	}

	if pay.Status() == string(dPayment.StatusSucceeded) {
		_ = s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusPaid)
	}

	return pay, nil
}

func (s *Service) HandleYooKassaWebhook(ctx context.Context, body []byte) error {
	var payload YooKassaWebhookPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}

	providerTxnID := strings.TrimSpace(payload.Object.ID)
	if providerTxnID == "" {
		return errors.New("empty yookassa payment id")
	}

	pay, err := s.repo.GetByProviderTxnID(ctx, providerTxnID)
	if err != nil {
		return err
	}

	switch payload.Event {
	case "payment.succeeded":
		pay.MarkSucceeded()

	case "payment.canceled":
		pay.MarkCanceled()

	default:
		s.applyPaymentStatus(pay, payload.Object.Status)
	}

	pay, err = s.repo.Update(ctx, pay)
	if err != nil {
		return err
	}

	if pay.Status() == string(dPayment.StatusSucceeded) {
		return s.orderRepo.UpdateStatus(ctx, pay.OrderID(), dOrder.OrderStatusPaid)
	}

	return nil
}

func (s *Service) applyPaymentStatus(pay *dPayment.Payment, status string) {
	switch normalizePaymentStatus(status) {
	case string(dPayment.StatusSucceeded):
		pay.MarkSucceeded()
	case string(dPayment.StatusCanceled):
		pay.MarkCanceled()
	case string(dPayment.StatusFailed):
		pay.MarkFailed()
	default:
		// pending оставляем как есть
	}
}

func normalizePaymentStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "succeeded", "success", "paid":
		return string(dPayment.StatusSucceeded)
	case "canceled", "cancelled":
		return string(dPayment.StatusCanceled)
	case "failed", "error":
		return string(dPayment.StatusFailed)
	default:
		return string(dPayment.StatusPending)
	}
}
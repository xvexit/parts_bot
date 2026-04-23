package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	paymentuc "partsBot/internal/usecase/payment"
	useruc "partsBot/internal/usecase/user"
)

type PaymentHandler struct {
	paymentService *paymentuc.Service
	userService    *useruc.Service
}

func NewPaymentHandler(paymentService *paymentuc.Service, userService *useruc.Service) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		userService:    userService,
	}
}

// CreatePayment — POST /api/user/orders/{id}/pay
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderIDStr := mux.Vars(r)["id"]
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil || orderID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	var req struct {
		ReturnURL string `json:"return_url"`
	}

	_ = json.NewDecoder(r.Body).Decode(&req)

	pay, err := h.paymentService.CreatePaymentForOrder(r.Context(), paymentuc.CreatePaymentInput{
		UserID:    userID,
		OrderID:   orderID,
		ReturnURL: req.ReturnURL,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, PaymentResponse{
		ID:            pay.ID(),
		OrderID:       pay.OrderID(),
		Amount:        pay.Amount(),
		Provider:      pay.Provider(),
		ProviderTxnID: pay.ProviderTxnID(),
		PaymentURL:    pay.PaymentURL(),
		Status:        pay.Status(),
		CreatedAt:     pay.CreatedAt(),
		UpdatedAt:     pay.UpdatedAt(),
	})
}

// GetLastPayment — GET /api/user/orders/{id}/payment
func (h *PaymentHandler) GetLastPayment(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderIDStr := mux.Vars(r)["id"]
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil || orderID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	pay, err := h.paymentService.GetLastByOrderID(r.Context(), userID, orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Payment not found")
		return
	}

	writeJSON(w, http.StatusOK, PaymentResponse{
		ID:            pay.ID(),
		OrderID:       pay.OrderID(),
		Amount:        pay.Amount(),
		Provider:      pay.Provider(),
		ProviderTxnID: pay.ProviderTxnID(),
		PaymentURL:    pay.PaymentURL(),
		Status:        pay.Status(),
		CreatedAt:     pay.CreatedAt(),
		UpdatedAt:     pay.UpdatedAt(),
	})
}

// ListPaymentsByOrder — GET /api/user/orders/{id}/payments
func (h *PaymentHandler) ListPaymentsByOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderIDStr := mux.Vars(r)["id"]
	orderID, err := strconv.ParseInt(strings.TrimSpace(orderIDStr), 10, 64)
	if err != nil || orderID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	payments, err := h.paymentService.ListByOrderID(r.Context(), userID, orderID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := make([]PaymentResponse, 0, len(payments))
	for _, pay := range payments {
		p := pay
		resp = append(resp, PaymentResponse{
			ID:            p.ID(),
			OrderID:       p.OrderID(),
			Amount:        p.Amount(),
			Provider:      p.Provider(),
			ProviderTxnID: p.ProviderTxnID(),
			PaymentURL:    p.PaymentURL(),
			Status:        p.Status(),
			CreatedAt:     p.CreatedAt(),
			UpdatedAt:     p.UpdatedAt(),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

// MockConfirm — POST /api/payments/mock/confirm
func (h *PaymentHandler) MockConfirm(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProviderTxnID string `json:"provider_txn_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	pay, err := h.paymentService.MarkSucceededByProviderTxnID(r.Context(), req.ProviderTxnID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, PaymentResponse{
		ID:            pay.ID(),
		OrderID:       pay.OrderID(),
		Amount:        pay.Amount(),
		Provider:      pay.Provider(),
		ProviderTxnID: pay.ProviderTxnID(),
		PaymentURL:    pay.PaymentURL(),
		Status:        pay.Status(),
		CreatedAt:     pay.CreatedAt(),
		UpdatedAt:     pay.UpdatedAt(),
	})
}
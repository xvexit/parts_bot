package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	paymentuc "partsBot/internal/usecase/payment"
	useruc "partsBot/internal/usecase/user"

	"github.com/gorilla/mux"
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

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderID, err := strconv.ParseInt(strings.TrimSpace(mux.Vars(r)["id"]), 10, 64)
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

	writeJSON(w, http.StatusCreated, paymentToResponse(pay))
}

func (h *PaymentHandler) GetLastPayment(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderID, err := strconv.ParseInt(strings.TrimSpace(mux.Vars(r)["id"]), 10, 64)
	if err != nil || orderID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	pay, err := h.paymentService.GetLastByOrderID(r.Context(), userID, orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Payment not found")
		return
	}

	writeJSON(w, http.StatusOK, paymentToResponse(pay))
}

func (h *PaymentHandler) ListPaymentsByOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderID, err := strconv.ParseInt(strings.TrimSpace(mux.Vars(r)["id"]), 10, 64)
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
	for i := range payments {
		resp = append(resp, paymentToResponse(&payments[i]))
	}

	writeJSON(w, http.StatusOK, resp)
}

// SyncOrderPayment — POST /api/user/orders/{id}/payment/sync
func (h *PaymentHandler) SyncOrderPayment(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderID, err := strconv.ParseInt(strings.TrimSpace(mux.Vars(r)["id"]), 10, 64)
	if err != nil || orderID <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order id")
		return
	}

	pay, err := h.paymentService.SyncOrderPayment(r.Context(), userID, orderID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, paymentToResponse(pay))
}

// YooKassaWebhook — POST /api/payments/yookassa/webhook
func (h *PaymentHandler) YooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read body")
		return
	}
	defer r.Body.Close()

	if err := h.paymentService.HandleYooKassaWebhook(r.Context(), body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func paymentToResponse(pay interface {
	ID() int64
	OrderID() int64
	Amount() int64
	Provider() string
	ProviderTxnID() *string
	PaymentURL() *string
	Status() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}) PaymentResponse {
	return PaymentResponse{
		ID:            pay.ID(),
		OrderID:       pay.OrderID(),
		Amount:        pay.Amount(),
		Provider:      pay.Provider(),
		ProviderTxnID: pay.ProviderTxnID(),
		PaymentURL:    pay.PaymentURL(),
		Status:        pay.Status(),
		CreatedAt:     pay.CreatedAt(),
		UpdatedAt:     pay.UpdatedAt(),
	}
}

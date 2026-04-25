package handler

import (
	"net/http"
	"strconv"

	orderuc "partsBot/internal/usecase/order"
	paymentuc "partsBot/internal/usecase/payment"
	useruc "partsBot/internal/usecase/user"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderService   *orderuc.Service
	paymentService *paymentuc.Service
	userService    *useruc.Service
}

func NewOrderHandler(
	orderService *orderuc.Service,
	paymentService *paymentuc.Service,
	userService *useruc.Service,
) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		paymentService: paymentService,
		userService:    userService,
	}
}

// ListOrders — GET /api/user/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orders, err := h.orderService.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	if len(orders) == 0 {
		writeJSON(w, http.StatusOK, "Список заказов пуст")
		return
	}

	resp := make([]OrderResponse, 0, len(orders))

	for _, ord := range orders {
		orderResp := toOrderResponse(ord)

		lastPayment, err := h.paymentService.GetLastByOrderID(r.Context(), userID, ord.ID())
		if err == nil && lastPayment != nil {
			if lastPayment.Status() == "pending" {
				if synced, syncErr := h.paymentService.SyncOrderPayment(r.Context(), userID, ord.ID()); syncErr == nil {
					lastPayment = synced
				}
			}

			orderResp.PaymentStatus = lastPayment.Status()
			orderResp.PaymentURL = lastPayment.PaymentURL()
		}

		resp = append(resp, orderResp)
	}

	writeJSON(w, http.StatusOK, resp)
}

// OrderItems — GET /api/user/orders/{order_id}/items
func (h *OrderHandler) OrderItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["order_id"]

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	items, err := h.orderService.OrderItems(r.Context(), orderID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get order items")
		return
	}

	writeJSON(w, http.StatusOK, items)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	orderuc "partsBot/internal/usecase/order"
	useruc "partsBot/internal/usecase/user"
)

type OrderHandler struct {
	orderService *orderuc.Service
	userService  *useruc.Service
}

func NewOrderHandler(orderService *orderuc.Service, userService *useruc.Service) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		userService:  userService,
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

	writeJSON(w, http.StatusOK, orders)
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
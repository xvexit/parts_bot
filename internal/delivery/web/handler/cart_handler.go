package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	cartuc "partsBot/internal/usecase/cart"
	useruc "partsBot/internal/usecase/user"

	"github.com/gorilla/mux"
)

type CartHandler struct {
	cartService *cartuc.Service
	userService *useruc.Service
}

func NewCartHandler(cartService *cartuc.Service, userService *useruc.Service) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		userService: userService,
	}
}

// ShowCart — GET /api/user/cart
func (h *CartHandler) ShowCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	items, err := h.cartService.GetCart(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	if len(items) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"items":   []any{},
			"total":   0,
			"message": "Корзина пуста",
		})
		return
	}

	var total int64 = 0
	responseItems := make([]map[string]any, len(items))

	for i, item := range items {
		itemTotal := item.Price().Amount() * item.Quantity()
		total += itemTotal

		responseItems[i] = map[string]any{
			"part_id":      item.PartID(),
			"name":         item.Name(),
			"brand":        item.Brand(),
			"price":        item.Price().Amount(),
			"quantity":     item.Quantity(),
			"delivery_day": item.DeliveryDay(),
			"total":        itemTotal,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": responseItems,
		"total": total,
	})
}

// AddItem — POST /api/user/cart/items
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CartItemDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	dto := cartuc.CartItemDto{
		PartID:      req.PartID,
		Name:        req.Name,
		Brand:       req.Brand,
		Price:       req.Price, //can test hardcode
		Quantity:    req.Quantity,
		DeliveryDay: req.DeliveryDay,
		ImageURL:    req.ImageURL,
	}

	if err := h.cartService.AddItem(r.Context(), userID, dto); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Товар успешно добавлен в корзину",
	})
}

// Checkout — POST /api/user/cart/checkout
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Address == "" {
		writeError(w, http.StatusBadRequest, "Address is required")
		return
	}

	order, err := h.cartService.Checkout(r.Context(), userID, req.Address)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"order_id": order.ID(),
		"message":  "Заказ успешно создан",
	})
}

// RemoveItem — DELETE /api/user/cart/items/{part_id}
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	partID := strings.TrimSpace(mux.Vars(r)["part_id"])
	if partID == "" {
		writeError(w, http.StatusBadRequest, "part_id is required")
		return
	}

	if err := h.cartService.RemoveItem(r.Context(), userID, partID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Товар удален из корзины",
	})
}

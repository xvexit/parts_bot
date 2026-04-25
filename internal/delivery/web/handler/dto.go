package handler

import (
	"partsBot/internal/domain/order"
	"time"
)

type UserDto struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CarDto struct {
	Name string `json:"name"`
	VIN  string `json:"vin"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CartItemDto struct {
	PartID      string `json:"part_id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Price       int64  `json:"price"`
	Quantity    int64  `json:"quantity"`
	DeliveryDay int    `json:"delivery_day"`
	ImageURL    string `json:"image_url,omitempty"`
}

type OrderResponse struct {
	ID            int64               `json:"id"`
	UserID        int64               `json:"user_id"`
	Address       string              `json:"address"`
	Items         []OrderItemResponse `json:"items"`
	Total         int64               `json:"total"`
	Status        string              `json:"status"`
	PaymentStatus string              `json:"payment_status,omitempty"`
	PaymentURL    *string             `json:"payment_url,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

type OrderItemResponse struct {
	PartID      string `json:"part_id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Price       int64  `json:"price"`
	Quantity    int64  `json:"quantity"`
	DeliveryDay int    `json:"delivery_day"`
}

type PaymentResponse struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	Amount        int64     `json:"amount"`
	Provider      string    `json:"provider"`
	ProviderTxnID *string   `json:"provider_txn_id,omitempty"`
	PaymentURL    *string   `json:"payment_url,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func toOrderResponse(o order.Order) OrderResponse {
	return OrderResponse{
		ID:        o.ID(),
		UserID:    o.UserID(),
		Items:     mapItems(o.Items()),
		Address:   o.Address(),
		Total:     o.Total().Amount(),
		Status:    o.Status(),
		CreatedAt: o.CreatedAt(),
	}
}

func mapItems(items []order.OrderItem) []OrderItemResponse {
	res := make([]OrderItemResponse, 0, len(items))

	for _, i := range items {
		res = append(res, OrderItemResponse{
			PartID:      i.PartID(),
			Name:        i.Name(),
			Brand:       i.Brand(),
			Price:       i.Price().Amount(),
			Quantity:    i.Quantity(),
			DeliveryDay: i.DeliveryDay(),
		})
	}

	return res
}

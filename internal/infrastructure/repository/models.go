package repository

import (
	"time"
)

type UserModel struct {
	ID        int64
	Name      string
	Email     *string
	Phone     string
	Password  string
	CreatedAt time.Time
}

type PaymentModel struct {
	ID            int64
	OrderID       int64
	Amount        int64
	Provider      string
	ProviderTxnID *string
	PaymentURL    *string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CartItemModel struct {
	PartID      string
	Name        string
	Brand       string
	Price       int64
	Quantity    int64
	DeliveryDay int
	ImageURL    string
}

type CarModel struct {
	ID     int64
	UserID int64
	Name   string
	VIN    string
}

type OrderModel struct {
	ID        int64
	UserID    int64
	Address   string
	Total     int64
	Status    string
	CreatedAt time.Time
}

type OrderItemModel struct {
	OrderID     int64
	PartID      string
	Name        string
	Brand       string
	Price       int64
	Quantity    int64
	DeliveryDay int
}

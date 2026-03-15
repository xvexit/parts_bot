package repository

import "time"

type CartItemDto struct {
	PartID      string
	Name        string
	Brand       string
	Price       int64
	Quantity    int64
	DeliveryDay int
	ImageURL    string
}

type CarDTO struct {
	ID     int64
	UserID int64
	Name   string
	VIN    string
}

type OrderDTO struct {
	ID        int64
	UserID    int64
	Address   string
	Status    string
	CreatedAt time.Time
}

type OrderItemDTO struct {
	OrderID     int64
	PartID      string
	Name        string
	Brand       string
	Price       int64
	Quantity    int64
	DeliveryDay int
}

type UserDTO struct {
	ID         int64
	TelegramID int64
	Name       string
	Phone      string
	CreatedAt  time.Time
}

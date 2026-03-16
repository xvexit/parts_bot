package telegram

import (
	"partsBot/internal/usecase/order"
)

type orderHandler struct {
	service *order.Service
}

func NewOrderHandler(service *order.Service) *orderHandler{
	return &orderHandler{
		service: service,
	}
}
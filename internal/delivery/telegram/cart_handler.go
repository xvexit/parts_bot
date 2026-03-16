package telegram

import "partsBot/internal/usecase/cart"

type cartHandler struct {
	service *cart.Service
}

func NewCartHandler(service *cart.Service) *cartHandler{
	return &cartHandler{
		service: service,
	}
}
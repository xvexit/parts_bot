package telegram

import "partsBot/internal/usecase/car"

type carHandler struct {
	service *car.Service
}

func NewCarHandler(service *car.Service) *carHandler{
	return &carHandler{
		service: service,
	}
}
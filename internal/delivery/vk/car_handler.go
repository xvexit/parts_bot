package vk

import (
	"context"
	"log"
	"strings"

	"partsBot/internal/usecase/car"
	"partsBot/internal/usecase/user"
)

type carHandler struct {
	service  *car.Service
	uService *user.Service
	bot      *Bot
}

func NewCarHandler(
	service *car.Service,
	uService *user.Service,
	bot *Bot,
) *carHandler {
	return &carHandler{
		service:  service,
		uService: uService,
		bot:      bot,
	}
}

func (h *carHandler) AddCar(userID int, text string) {

	ctx := context.Background()
	
	args := strings.Fields(text)
	if len(args) < 2 {
		h.bot.sendMessage(userID, "Ошибка! Используй: addcar марка VIN")
		return
	}

	us, err := h.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		h.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	dto := car.CarDto{
		UserID: us.ID(),
		Name:   args[0],
		VIN:    args[1],
	}

	ccar, err := h.service.Add(ctx, dto)
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка добавления авто!")
		log.Println(err)
		return
	}

	msg := "Авто " + ccar.Name() + " с VIN " + ccar.Vin() + " добавлено!"
	h.bot.sendMessage(userID, msg)
}

func (h *carHandler) MyCars(userID int) {

	ctx := context.Background()

	us, err := h.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		h.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	cars, err := h.service.ListByUser(ctx, us.ID())
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка получения списка авто!")
		log.Println(err)
		return
	}

	if len(cars) == 0 {
		h.bot.sendMessage(userID, "У вас нет добавленных авто")
		return
	}

	result := "Ваши авто:\n"

	for _, c := range cars {
		result += c.Name() + " VIN: " + c.Vin() + "\n"
	}

	h.bot.sendMessage(userID, result)
}

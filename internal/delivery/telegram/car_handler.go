package telegram

import (
	"context"
	"log"
	"partsBot/internal/usecase/car"
	"partsBot/internal/usecase/user"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type carHandler struct {
	service  *car.Service
	uService *user.Service
}

func NewCarHandler(service *car.Service, uService *user.Service) *carHandler {
	return &carHandler{
		service:  service,
		uService: uService,
	}
}

func (h *carHandler) AddCar(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	ctx := context.Background()

	args := strings.Split(msg.CommandArguments(), " ")

	if len(args) < 2 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка добавления авто! Пишите /addcar марка вин. Ровно так, через пробел"))
		return
	}

	us, err := h.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	dto := car.CarDto{
		UserID: us.ID(),
		Name:   args[0],
		VIN:    args[1],
	}

	c, err := h.service.Add(ctx, dto)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка добавления авто! Пишите /addcar марка вин. Ровно так, через пробел"))
		log.Println(err)
		return
	}

	text := "Авто " + c.Name() + " с вин номером " + c.Vin() + " успешно добавлен!"
	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func (h *carHandler) MyCars(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	ctx := context.Background()

	us, err := h.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	c, err := h.service.ListByUser(ctx, us.ID())
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка поиска авто!"))
		log.Println(err)
		return
	}

	if len(c) == 0{
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Вы не добавили ни одного авто!"))
		return
	}

	text := "Ваши авто:\n"
	for _, a := range c {
		text += a.Name() + " VIN: " + a.Vin() + "\n"
	}
	
	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

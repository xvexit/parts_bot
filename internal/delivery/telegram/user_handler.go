package telegram

import (
	"context"
	"log"
	"partsBot/internal/usecase/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *userHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) Start(api *tgbotapi.BotAPI, msg *tgbotapi.Message) { //надо получить phone
	ctx := context.Background()

	dto := user.UserDto{
		TelegramID: int64(msg.From.ID),
		Name:       msg.From.UserName,
		Phone:      "0",
	}

	u, err := h.service.Register(ctx, dto)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка регистрации"))
		log.Printf("ERROR: Регистрация не удалась: %v", err)
		return
	}

	text := u.Name() + ", вы успешно зарегестрированы!"

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}


package vk

import (
	"context"
	"log"

	"partsBot/internal/usecase/user"

	"github.com/SevereCloud/vksdk/v2/api"
)

type userHandler struct {
	service *user.Service
	bot     *Bot
}

func NewUserHandler(service *user.Service, bot *Bot) *userHandler {
	return &userHandler{
		service: service,
		bot:     bot,
	}
}

func (h *userHandler) Start(userID int) {

	ctx := context.Background()

	name, err := h.bot.GetUser(userID)
	if err != nil || name == "" {
		name = "Пользователь"
	}

	dto := user.UserDto{
		TelegramID: int64(userID), // переиспользуешь поле как external_id
		Name:       name,          // VK не всегда даёт username
		Phone:      "0",           // временно
	}

	log.Println(dto.Name)

	u, err := h.service.Register(ctx, dto)
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка регистрации")
		log.Printf("ERROR: регистрация: %v", err)
		return
	}

	msg := u.Name() + ", вы успешно зарегистрированы!"
	h.bot.sendMessage(userID, msg)
}

func (b *Bot) GetUser(userID int) (string, error) {
	res, err := b.vk.UsersGet(api.Params{
		"user_ids": userID,
		"fields":   "first_name,last_name",
	})
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", nil
	}

	u := res[0]

	return u.FirstName + " " + u.LastName, nil
}

package vk

import (
	"context"
	"fmt"
	"log"
	"strings"

	adap "partsBot/internal/infrastructure/adapter"
	"partsBot/internal/usecase/car"
	"partsBot/internal/usecase/user"
)

type partsHandler struct {
	parts       *adap.Adapter
	carService  *car.Service
	userService *user.Service
	bot         *Bot
}

func NewPartsHandler(
	parts *adap.Adapter,
	carService *car.Service,
	userService *user.Service,
	bot *Bot,
) *partsHandler {
	return &partsHandler{
		parts:       parts,
		carService:  carService,
		userService: userService,
		bot:         bot,
	}
}

func (h *partsHandler) Search(userID int, text string) {

	ctx := context.Background()

	us, err := h.userService.GetByTgID(ctx, int64(userID))
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка! Пользователь не найден!")
		log.Println(err)
		return
	}

	cars, err := h.carService.ListByUser(ctx, us.ID())
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка! Авто не найдены!")
		log.Println(err)
		return
	}

	if len(cars) == 0 {
		h.bot.sendMessage(userID, "У вас нет добавленных авто")
		return
	}

	searchText := strings.TrimSpace(text)
	if searchText == "" {
		h.bot.sendMessage(userID, "Введите запрос для поиска")
		return
	}

	vin := cars[0].Vin()

	parts, err := h.parts.GetPartsByVIN(ctx, vin, searchText, false)
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка поиска товаров!")
		log.Println(err)
		return
	}

	if len(parts) == 0 {
		h.bot.sendMessage(userID, "По запросу ничего не найдено")
		return
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"Найдено групп: %d\nПо запросу: %s\n\n",
		len(parts),
		searchText,
	))

	q := strings.ToLower(searchText)
	found := 0

	for _, g := range parts {
		if strings.Contains(strings.ToLower(g.Name+" "+g.Shortname), q) {
			found++
			sb.WriteString(fmt.Sprintf(
				"• %s (%s)\n  Артикулы: %s\n\n",
				g.Name,
				g.Shortname,
				g.Parts,
			))
		}
	}

	if found == 0 {
		h.bot.sendMessage(userID, "По запросу ничего не найдено")
		return
	}

	h.bot.sendMessage(userID, sb.String())
}

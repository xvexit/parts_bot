package telegram

import (
	"context"
	"fmt"
	"log"
	adap "partsBot/internal/infrastructure/adapter"
	"partsBot/internal/usecase/car"
	"partsBot/internal/usecase/user"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// parts_handler.go
type partsHandler struct {
	parts       *adap.Adapter
	carService  *car.Service
	userService *user.Service
}

func NewPartsHandler(parts *adap.Adapter, carService *car.Service, userService *user.Service) *partsHandler {
	return &partsHandler{
		parts:       parts,
		carService:  carService,
		userService: userService,
	}
}

func (h *partsHandler) Search(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	us, err := h.userService.GetByTgID(ctx, msg.From.ID)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка! Пользователь не найден!"))
	}

	cars, err := h.carService.ListByUser(ctx, us.ID())
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка! Авто не были найдены!"))
	}

	searchText := msg.CommandArguments()

	parts, err := h.parts.GetPartsByVIN(ctx, cars[0].Vin(), searchText, false)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка поиска товара!"))
		log.Println(err)
		return
	}

	if len(parts) == 0 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "По этому VIN ничего не найдено"))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Найдено групп: %d\nПо запросу: %s\n\n", len(parts), searchText))

	found := 0
	q := strings.ToLower(searchText)

	for _, g := range parts {
		if strings.Contains(strings.ToLower(g.Name+" "+g.Shortname), q) {
			found++
			sb.WriteString(fmt.Sprintf("• %s (%s)\n  Артикулы: %s\n\n", g.Name, g.Shortname, g.Parts))
		}
	}

	if found == 0 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "По запросу «"+searchText+"» ничего не найдено."))
		return
	}

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, sb.String()))
}

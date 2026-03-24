package telegram

import (
	"context"
	"fmt"
	"log"
	"partsBot/internal/usecase/order"
	"partsBot/internal/usecase/user"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type orderHandler struct {
	service  *order.Service
	uService *user.Service
}

func NewOrderHandler(service *order.Service, uService *user.Service) *orderHandler {
	return &orderHandler{
		service: service,
		uService: uService,
	}
}

func (h *orderHandler) ListOrders(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	us, err := h.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil{
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	orders, err := h.service.ListByUserID(ctx, us.ID())
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка поиска списка заказов!"))
		log.Println(err)
		return
	}

	if len(orders) == 0 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Список заказов пуст!"))
		return
	}

	text := "Заказы:\n"

	for _, o := range orders {
		text += fmt.Sprintf("Заказ %d на адрес %s Статус заказа %s\n", o.ID(), o.Address(), o.Status())
	}

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func (h *orderHandler) OrderItems(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	id, err := strconv.ParseInt(msg.CommandArguments(), 10, 64)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Введеный номер заказа не валидный!"))
		log.Println(err)
		return
	}

	items, err := h.service.OrderItems(ctx, id)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка поиска товаров!"))
		log.Println(err)
		return
	}

	text := "Товары:\n"

	for _, i := range items {
		text += fmt.Sprintf("%s %s %d x%d\n", i.Brand(), i.Name(), i.Price().Amount(), i.Quantity())
	}
	api.Send(tgbotapi.NewMessage(msg.Chat.ID, "text"))
}

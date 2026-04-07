package vk

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"partsBot/internal/usecase/order"
	"partsBot/internal/usecase/user"
)

type orderHandler struct {
	service  *order.Service
	uService *user.Service
	bot      *Bot
}

func NewOrderHandler(
	service *order.Service,
	uService *user.Service,
	bot *Bot,
) *orderHandler {
	return &orderHandler{
		service:  service,
		uService: uService,
		bot:      bot,
	}
}

func (h *orderHandler) ListOrders(userID int) {

	ctx := context.Background()

	us, err := h.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		h.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	orders, err := h.service.ListByUserID(ctx, us.ID())
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка поиска заказов!")
		log.Println(err)
		return
	}

	if len(orders) == 0 {
		h.bot.sendMessage(userID, "Список заказов пуст!")
		return
	}

	result := "Заказы:\n"

	for _, o := range orders {
		result += fmt.Sprintf(
			"Заказ %d | адрес: %s | статус: %s\n",
			o.ID(),
			o.Address(),
			o.Status(),
		)
	}

	h.bot.sendMessage(userID, result)
}

func (h *orderHandler) OrderItems(userID int, text string) {

	ctx := context.Background()

	id, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		h.bot.sendMessage(userID, "Номер заказа введён неверно!")
		log.Println(err)
		return
	}

	items, err := h.service.OrderItems(ctx, id)
	if err != nil {
		h.bot.sendMessage(userID, "Ошибка получения товаров заказа!")
		log.Println(err)
		return
	}

	if len(items) == 0 {
		h.bot.sendMessage(userID, "В заказе нет товаров!")
		return
	}

	result := "Товары:\n"

	for _, i := range items {
		result += fmt.Sprintf(
			"%s %s %d x%d\n",
			i.Brand(),
			i.Name(),
			i.Price().Amount(),
			i.Quantity(),
		)
	}

	h.bot.sendMessage(userID, result)
}

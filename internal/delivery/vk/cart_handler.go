package vk

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"partsBot/internal/usecase/cart"
	"partsBot/internal/usecase/user"
)

type cartHandler struct {
	service  *cart.Service
	uService *user.Service
	bot      *Bot
}

func NewCartHandler(
	service *cart.Service,
	uService *user.Service,
	bot *Bot,
) *cartHandler {
	return &cartHandler{
		service:  service,
		uService: uService,
		bot:      bot,
	}
}

func (c *cartHandler) ShowCart(userID int) {

	ctx := context.Background()

	us, err := c.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		c.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	items, err := c.service.GetCart(ctx, us.ID())
	if err != nil {
		c.bot.sendMessage(userID, "Ошибка поиска корзины!")
		log.Println(err)
		return
	}

	if len(items) == 0 {
		c.bot.sendMessage(userID, "Корзина пуста!")
		return
	}

	result := "Ваша корзина:\n"

	for _, item := range items {
		result += fmt.Sprintf(
			"%s %s x%d\n",
			item.Brand(),
			item.Name(),
			item.Quantity(),
		)
	}

	c.bot.sendMessage(userID, result)
}

func (c *cartHandler) AddItem(userID int, text string) {

	ctx := context.Background()

	args := strings.Fields(text)
	if len(args) < 6 {
		c.bot.sendMessage(userID, "Неверные аргументы. Пример: additem partID name brand price qty delivDay")
		return
	}

	qty, err := strconv.Atoi(args[4])
	if err != nil {
		c.bot.sendMessage(userID, "Количество должно быть числом")
		return
	}

	dd, err := strconv.Atoi(args[5])
	if err != nil {
		c.bot.sendMessage(userID, "Дни доставки должны быть числом")
		return
	}

	dto := cart.CartItemDto{
		PartID:      args[0],
		Name:        args[1],
		Brand:       args[2],
		Quantity:    int64(qty),
		DeliveryDay: dd,
		ImageURL:    "some",
	}

	us, err := c.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		c.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	err = c.service.AddItem(ctx, us.ID(), dto)
	if err != nil {
		c.bot.sendMessage(userID, "Ошибка добавления товара!")
		log.Println(err)
		return
	}

	c.bot.sendMessage(userID, "Товар добавлен в корзину!")
}

func (c *cartHandler) Checkout(userID int, text string) {

	ctx := context.Background()

	address := text

	us, err := c.uService.GetByTgID(ctx, int64(userID))
	if err != nil {
		c.bot.sendMessage(userID, "Пользователь не найден!")
		log.Println(err)
		return
	}

	ord, err := c.service.Checkout(ctx, us.ID(), address)
	if err != nil {
		c.bot.sendMessage(userID, "Ошибка создания заказа!")
		log.Println(err)
		return
	}

	msg := fmt.Sprintf("Заказ создан! ID: %d", ord.ID())

	c.bot.sendMessage(userID, msg)
}

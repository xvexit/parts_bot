package telegram

import (
	"context"
	"fmt"
	"log"
	"partsBot/internal/usecase/cart"
	"partsBot/internal/usecase/user"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type cartHandler struct {
	service *cart.Service
	uService *user.Service
}

func NewCartHandler(service *cart.Service, uService *user.Service) *cartHandler {
	return &cartHandler{
		service: service,
		uService: uService,
	}
}

func (c *cartHandler) ShowCart(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	us, err := c.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil{
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	items, err := c.service.GetCart(ctx, us.ID())
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка поиска корзины!"))
		log.Println(err)
		return
	}

	if len(items) == 0 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Корзина пуста!"))
		return
	}

	text := "Ваша корзина:\n"

	for _, item := range items {
		text += fmt.Sprintf(
			"%s %s x%d\n",
			item.Brand(),
			item.Name(),
			item.Quantity(),
		)
	}

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func (c *cartHandler) AddItem(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	args := strings.Split(msg.CommandArguments(), " ")
	if len(args) < 6 {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неверное количество аргументов. /additem partID name brand price qty delivDay"))
		log.Println(args)
		return
	}

	qty, err := strconv.Atoi(args[3])
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Количество должно быть числом"))
		return
	}

	dd, err := strconv.Atoi(args[3])
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Количество дней доставки должно быть числом"))
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

	us, err := c.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil{
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	if err := c.service.AddItem(ctx, us.ID(), dto); err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка добавления товара в корзину!"))
		log.Println(err)
		return
	}

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Товар успешно добавлен в корзину!"))
}

func (c *cartHandler) Checkout(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	ctx := context.Background()

	address := msg.CommandArguments()

	us, err := c.uService.GetByTgID(ctx, msg.From.ID)
	if err != nil{
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пользователь не найден!"))
		log.Println(err)
		return
	}

	order, err := c.service.Checkout(ctx, us.ID(), address)
	if err != nil {
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка создания заказа"))
		log.Println(err)
		return
	}

	text := fmt.Sprintf("Заказ создан! Номер заказа: ", order.ID())

	api.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

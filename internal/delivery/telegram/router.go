package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	userHandler  *userHandler
	carHandler   *carHandler
	cartHandler  *cartHandler
	orderHandler *orderHandler
	partsHandler *partsHandler
}

func NewRouter(
	userHandler *userHandler,
	carHandler *carHandler,
	cartHandler *cartHandler,
	orderHandler *orderHandler,
	partsHandler *partsHandler,
) *Router {
	return &Router{
		userHandler:  userHandler,
		carHandler:   carHandler,
		cartHandler:  cartHandler,
		orderHandler: orderHandler,
		partsHandler: partsHandler,
	}
}

func (r *Router) Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Printf("Входящее сообщение: [%s] %s", msg.From.UserName, msg.Text)
	switch msg.Command() {
	case "start":
		r.userHandler.Start(api, msg)
	case "cart":
		r.cartHandler.ShowCart(api, msg)
	case "additem":
		r.cartHandler.AddItem(api, msg)
	case "checkout":
		r.cartHandler.Checkout(api, msg)
	case "addcar":
		r.carHandler.AddCar(api, msg)
	case "cars":
		r.carHandler.MyCars(api, msg)
	case "orders":
		r.orderHandler.ListOrders(api, msg)
	case "orderitems":
		r.orderHandler.OrderItems(api, msg)
	case "getparts":
		r.partsHandler.Search(api, msg)
	default:
		api.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная комманда!"))
	}
}

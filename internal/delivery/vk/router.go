package vk

import (
	"log"
	"strings"
)

type Router struct {
	userHandler  *userHandler
	carHandler   *carHandler
	cartHandler  *cartHandler
	orderHandler *orderHandler
	partsHandler *partsHandler
	bot          *Bot
}

func NewRouter(
	userHandler *userHandler,
	carHandler *carHandler,
	cartHandler *cartHandler,
	orderHandler *orderHandler,
	partsHandler *partsHandler,
	bot *Bot,
) *Router {
	return &Router{
		userHandler:  userHandler,
		carHandler:   carHandler,
		cartHandler:  cartHandler,
		orderHandler: orderHandler,
		partsHandler: partsHandler,
		bot:          bot,
	}
}

func (r *Router) Handle(userID int, text string) {

	log.Printf("VK IN [%d]: %s", userID, text)

	text = strings.TrimSpace(text)

	switch {

	// USER
	case text == "/start":
		r.userHandler.Start(userID)

	// CART
	case text == "/cart":
		r.cartHandler.ShowCart(userID)

	case strings.HasPrefix(text, "/additem"):
		r.cartHandler.AddItem(userID, strings.TrimPrefix(text, "/additem"))

	case strings.HasPrefix(text, "/checkout"):
		r.cartHandler.Checkout(userID, strings.TrimPrefix(text, "/checkout"))

	// CAR
	case strings.HasPrefix(text, "/addcar"):
		r.carHandler.AddCar(userID, strings.TrimPrefix(text, "/addcar"))

	case text == "/cars":
		r.carHandler.MyCars(userID)

	// ORDERS
	case text == "/orders":
		r.orderHandler.ListOrders(userID)

	case strings.HasPrefix(text, "/orderitems"):
		r.orderHandler.OrderItems(userID, strings.TrimPrefix(text, "/orderitems"))

	// PARTS
	case strings.HasPrefix(text, "/getparts"):
		r.partsHandler.Search(userID, strings.TrimPrefix(text, "/getparts"))

	default:
		r.sendUnknown(userID)
	}
}

func (r *Router) sendUnknown(userID int) {
	r.bot.sendMessage(userID, "Неизвестная команда. Напишите /start")
}

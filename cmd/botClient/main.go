package main

import (
	"fmt"
	"log"
	"os"

	"partsBot/internal/delivery/telegram"
	"partsBot/internal/infrastructure/db"

	repo "partsBot/internal/infrastructure/repository"

	caruc "partsBot/internal/usecase/car"
	cartuc "partsBot/internal/usecase/cart"
	orderuc "partsBot/internal/usecase/order"
	useruc "partsBot/internal/usecase/user"
)

func main() {
	dbConn, err := db.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repo.NewPostgresUserRepo(dbConn)
	carRepo := repo.NewPostgresCarRepo(dbConn)
	cartRepo := repo.NewPostgresCartRepo(dbConn)
	orderRepo := repo.NewPostgresOrderRepo(dbConn)

	txManager := db.NewTxManager(dbConn.Pool())

	userService := useruc.NewService(userRepo)
	carService := caruc.NewService(carRepo)
	cartService := cartuc.NewService(cartRepo, orderRepo, txManager)
	orderService := orderuc.NewService(orderRepo)

	fmt.Println(carService, cartService, orderService)

	userHandler := telegram.NewUserHandler(userService)
	cartHandler := telegram.NewCartHandler(cartService)
	carHandler := telegram.NewCarHandler(carService)
	orderHandler := telegram.NewOrderHandler(orderService)

	router := telegram.NewRouter(
		userHandler,
		carHandler,
		cartHandler,
		orderHandler,
	)

	bot, err := telegram.NewBot(os.Getenv("TELEGRAM_TOKEN"), router)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}

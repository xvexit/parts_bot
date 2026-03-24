package main

import (
	"fmt"
	"os"
	"log"

	"partsBot/internal/delivery/telegram"
	"partsBot/internal/infrastructure/db"

	repo "partsBot/internal/infrastructure/repository"

	caruc "partsBot/internal/usecase/car"
	cartuc "partsBot/internal/usecase/cart"
	orderuc "partsBot/internal/usecase/order"
	useruc "partsBot/internal/usecase/user"
	adapter "partsBot/internal/infrastructure/adapter"
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
	partAdapter := adapter.New(adapter.Config{
		APIKey: os.Getenv("PART_API"),
	})

	fmt.Println(carService, cartService, orderService)

	
	userHandler := telegram.NewUserHandler(userService)
	cartHandler := telegram.NewCartHandler(cartService, userService)
	carHandler := telegram.NewCarHandler(carService, userService)
	orderHandler := telegram.NewOrderHandler(orderService, userService)
	partHandler := telegram.NewPartsHandler(partAdapter, carService, userService)

	router := telegram.NewRouter(
		userHandler,
		carHandler,
		cartHandler,
		orderHandler,
		partHandler,
	)

	bot, err := telegram.NewBot(os.Getenv("TELEGRAM_TOKEN"), router)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}

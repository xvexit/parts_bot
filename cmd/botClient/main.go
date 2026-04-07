package main

import (
	"log"
	"os"
	"strconv"

	"partsBot/internal/infrastructure/db"
	repo "partsBot/internal/infrastructure/repository"

	caruc "partsBot/internal/usecase/car"
	cartuc "partsBot/internal/usecase/cart"
	orderuc "partsBot/internal/usecase/order"
	useruc "partsBot/internal/usecase/user"

	adap "partsBot/internal/infrastructure/adapter"

	vk "partsBot/internal/delivery/vk"
)

func main() {

	// DB
	dbConn, err := db.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Repos
	userRepo := repo.NewPostgresUserRepo(dbConn)
	carRepo := repo.NewPostgresCarRepo(dbConn)
	cartRepo := repo.NewPostgresCartRepo(dbConn)
	orderRepo := repo.NewPostgresOrderRepo(dbConn)

	txManager := db.NewTxManager(dbConn.Pool())

	// Services
	userService := useruc.NewService(userRepo)
	carService := caruc.NewService(carRepo)
	cartService := cartuc.NewService(cartRepo, orderRepo, txManager)
	orderService := orderuc.NewService(orderRepo)

	partAdapter := adap.New(adap.Config{
		APIKey: os.Getenv("PART_API"),
	})

	// BOT CORE
	groupID, _ := strconv.Atoi(os.Getenv("GROUP_ID"))

	bot, err := vk.NewBot(os.Getenv("VK_API"), groupID)
	if err != nil {
		log.Fatal(err)
	}
	userHandler := vk.NewUserHandler(userService, bot)
	cartHandler := vk.NewCartHandler(cartService, userService, bot)
	carHandler := vk.NewCarHandler(carService, userService, bot)
	orderHandler := vk.NewOrderHandler(orderService, userService, bot)
	partHandler := vk.NewPartsHandler(partAdapter, carService, userService, bot)

	router := vk.NewRouter(
		userHandler,
		carHandler,
		cartHandler,
		orderHandler,
		partHandler,
		bot,
	)

	bot.SetRouter(router)

	bot.Start()
}

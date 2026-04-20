package main

import (
	"log"
	"os"

	"partsBot/internal/delivery/web"
	"partsBot/internal/delivery/web/handler"
	"partsBot/internal/infrastructure/db"
	repo "partsBot/internal/infrastructure/repository"

	caruc "partsBot/internal/usecase/car"
	cartuc "partsBot/internal/usecase/cart"
	orderuc "partsBot/internal/usecase/order"
	useruc "partsBot/internal/usecase/user"

	adap "partsBot/internal/infrastructure/adapter"
	jauth "partsBot/internal/infrastructure/auth"
	uauth "partsBot/internal/usecase/auth"
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
	tokenManager := jauth.NewJWTManager()
	userService := useruc.NewService(userRepo)
	carService := caruc.NewService(carRepo)
	cartService := cartuc.NewService(cartRepo, orderRepo, txManager)
	orderService := orderuc.NewService(orderRepo)
	authService := uauth.NewService(tokenManager, userRepo)

	partAdapter := adap.New(adap.Config{
		APIKey:     os.Getenv("PART_API"),
		APIKeyVIN:  os.Getenv("PART_API_VIN"),
		APIKeyTree: os.Getenv("PART_API_TREE"),
	})

	userHandler := handler.NewUserHandler(userService)
	cartHandler := handler.NewCartHandler(cartService, userService)
	carHandler := handler.NewCarHandler(carService, userService)
	orderHandler := handler.NewOrderHandler(orderService, userService)
	partHandler := handler.NewPartsHandler(partAdapter, carService, userService)
	authHandler := handler.NewAuthHandler(authService)

	router := web.NewRouter(
		authHandler,
		userHandler,
		carHandler,
		cartHandler,
		orderHandler,
		partHandler,
		tokenManager,
	)

	server := web.NewServer(router)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

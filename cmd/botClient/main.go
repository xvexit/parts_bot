package main

import (
	"log"
	"os"

	"partsBot/internal/delivery/web"
	"partsBot/internal/delivery/web/handler"
	adap "partsBot/internal/infrastructure/adapter"
	jauth "partsBot/internal/infrastructure/auth"
	"partsBot/internal/infrastructure/db"
	repo "partsBot/internal/infrastructure/repository"

	uauth "partsBot/internal/usecase/auth"
	caruc "partsBot/internal/usecase/car"
	cartuc "partsBot/internal/usecase/cart"
	orderuc "partsBot/internal/usecase/order"
	partuc "partsBot/internal/usecase/part"
	paymentuc "partsBot/internal/usecase/payment"
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
	partRepo := repo.NewPostgresPartCatalogRepo(dbConn)
	paymentRepo := repo.NewPostgresPaymentRepo(dbConn)

	txManager := db.NewTxManager(dbConn.Pool())

	tokenManager := jauth.NewJWTManager()
	userService := useruc.NewService(userRepo)
	carService := caruc.NewService(carRepo)
	cartService := cartuc.NewService(cartRepo, orderRepo, txManager)
	orderService := orderuc.NewService(orderRepo)
	partService := partuc.NewService(partRepo)
	authService := uauth.NewService(tokenManager, userRepo)

	paymentGateway := paymentuc.NewYooKassaGateway(
		getEnv("YOOKASSA_SHOP_ID", ""),
		getEnv("YOOKASSA_SECRET_KEY", ""),
	)
	paymentService := paymentuc.NewService(paymentRepo, orderRepo, paymentGateway)

	partAdapter := adap.New(adap.Config{
		APIKey:     os.Getenv("PART_API"),
		APIKeyVIN:  os.Getenv("PART_API_VIN"),
		APIKeyTree: os.Getenv("PART_API_TREE"),
	})

	userHandler := handler.NewUserHandler(userService)
	cartHandler := handler.NewCartHandler(cartService, userService)
	carHandler := handler.NewCarHandler(carService, userService)
	paymentHandler := handler.NewPaymentHandler(paymentService, userService)
	orderHandler := handler.NewOrderHandler(orderService, paymentService, userService)
	partHandler := handler.NewPartsHandler(partAdapter, carService, partService, userService)
	authHandler := handler.NewAuthHandler(authService)

	router := web.NewRouter(
		authHandler,
		userHandler,
		carHandler,
		cartHandler,
		orderHandler,
		paymentHandler,
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
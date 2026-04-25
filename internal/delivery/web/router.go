	package web

	import (
		"net/http"

		"partsBot/internal/delivery/web/handler"
		"partsBot/internal/delivery/web/handler/middleware"
		"partsBot/internal/usecase/auth"

		"github.com/gorilla/mux"
	)

	type Router struct {
		mux *mux.Router
	}

	func NewRouter(
		authH *handler.AuthHandler,
		userH *handler.UserHandler,
		carH *handler.CarHandler,
		cartH *handler.CartHandler,
		orderH *handler.OrderHandler,
		paymentH *handler.PaymentHandler,
		partsH *handler.PartsHandler,
		tm auth.TokenManager,
	) *Router {
		r := mux.NewRouter()

		api := r.PathPrefix("/api").Subrouter()

		api.HandleFunc("/auth/register", userH.Register).Methods("POST")
		api.HandleFunc("/auth/login", authH.Login).Methods("POST")
		api.HandleFunc("/auth/refresh", authH.Refresh).Methods("POST")
		api.HandleFunc("/parts/search", partsH.SearchParts).Methods("GET")
		api.HandleFunc("/payments/yookassa/webhook", paymentH.YooKassaWebhook).Methods("POST")

		userRouter := api.PathPrefix("/user").Subrouter()
		userRouter.Use(middleware.Auth(tm))

		userRouter.HandleFunc("/parts/tree", partsH.SearchTree).Methods("GET")
		userRouter.HandleFunc("/orders/{id:[0-9]+}/payment/sync", paymentH.SyncOrderPayment).Methods("POST")
		userRouter.HandleFunc("/parts/check", partsH.CheckPartOffer).Methods("POST")

		userRouter.HandleFunc("/cars", carH.MyCars).Methods("GET")
		userRouter.HandleFunc("/cars", carH.AddCar).Methods("POST")
		userRouter.HandleFunc("/cars/{id:[0-9]+}", carH.DeleteCar).Methods("DELETE")

		userRouter.HandleFunc("/cart", cartH.ShowCart).Methods("GET")
		userRouter.HandleFunc("/cart/items", cartH.AddItem).Methods("POST")
		userRouter.HandleFunc("/cart/items/{part_id}", cartH.RemoveItem).Methods("DELETE")
		userRouter.HandleFunc("/cart/checkout", cartH.Checkout).Methods("POST")

		userRouter.HandleFunc("/orders", orderH.ListOrders).Methods("GET")
		userRouter.HandleFunc("/orders/{order_id:[0-9]+}/items", orderH.OrderItems).Methods("GET")

		userRouter.HandleFunc("/orders/{id:[0-9]+}/pay", paymentH.CreatePayment).Methods("POST")
		userRouter.HandleFunc("/orders/{id:[0-9]+}/payment", paymentH.GetLastPayment).Methods("GET")
		userRouter.HandleFunc("/orders/{id:[0-9]+}/payments", paymentH.ListPaymentsByOrder).Methods("GET")

		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./internal/delivery/web/index.html")
		}).Methods("GET")

		r.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./internal/delivery/web/styles.css")
		}).Methods("GET")

		r.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./internal/delivery/web/main.js")
		}).Methods("GET")

		r.HandleFunc("/payment-result.html", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./internal/delivery/web/payment-result.html")
		}).Methods("GET")

		return &Router{mux: r}
	}

	func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		r.mux.ServeHTTP(w, req)
	}
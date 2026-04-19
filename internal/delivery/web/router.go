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
	partsH *handler.PartsHandler,
	tm auth.TokenManager,
) *Router {

	r := mux.NewRouter()

	// ========== API PREFIX ==========
	api := r.PathPrefix("/api").Subrouter()

	// PUBLIC
	api.HandleFunc("/auth/register", userH.Register).Methods("POST")
	api.HandleFunc("/auth/login", authH.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authH.Refresh).Methods("POST")
	api.HandleFunc("/parts/search", partsH.SearchParts).Methods("GET")

	// PROTECTED
	userRouter := api.PathPrefix("/user").Subrouter()
	userRouter.Use(middleware.Auth(tm))

	userRouter.HandleFunc("/cars", carH.MyCars).Methods("GET")
	userRouter.HandleFunc("/cars", carH.AddCar).Methods("POST")

	userRouter.HandleFunc("/cart", cartH.ShowCart).Methods("GET")
	userRouter.HandleFunc("/cart/items", cartH.AddItem).Methods("POST")
	userRouter.HandleFunc("/cart/checkout", cartH.Checkout).Methods("POST")

	userRouter.HandleFunc("/orders", orderH.ListOrders).Methods("GET")

	// ========== STATIC FRONTEND ==========
	// index.html uses relative links that browsers resolve to /styles.css and /main.js
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

// ServeHTTP нужен для совместимости с http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

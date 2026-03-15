package order

type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "new"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCanceled  OrderStatus = "canceled"
	OrderStatusErr       OrderStatus = "error"
)

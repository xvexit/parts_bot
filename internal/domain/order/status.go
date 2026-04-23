package order

type OrderStatus string

const (
	OrderStatusNew            OrderStatus = "new"
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusConfirmed      OrderStatus = "confirmed"
	OrderStatusDelivered      OrderStatus = "delivered"
	OrderStatusCanceled       OrderStatus = "canceled"
	OrderStatusErr            OrderStatus = "error"
)
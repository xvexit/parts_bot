package order

type OrderStatus string

const (
    OrderStatusNew       OrderStatus = "new"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCanceled  OrderStatus = "canceled"
)
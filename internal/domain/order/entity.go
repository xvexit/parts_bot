package order

import (
	money "partsBot/internal/domain/shared"
	"partsBot/pkg/errors"
	"strings"
	"time"
)

type Order struct {
	id          int64
	userID      int64
	address     string
	items       []OrderItem
	total       money.Money
	status      OrderStatus
	createdAt   time.Time
}

type OrderItem struct {
	partID      string
	name        string
	brand       string
	price       money.Money
	quantity    int64
	deliveryDay int
}

func RestoreOrder(
	id, userID, total int64,
	items []OrderItem,
	status, address string,
	createdAt time.Time,
) *Order {
	var st OrderStatus

	switch status {
	case string(OrderStatusNew):
		st = OrderStatusNew
	case string(OrderStatusPendingPayment):
		st = OrderStatusPendingPayment
	case string(OrderStatusPaid):
		st = OrderStatusPaid
	case string(OrderStatusConfirmed):
		st = OrderStatusConfirmed
	case string(OrderStatusDelivered):
		st = OrderStatusDelivered
	case string(OrderStatusCanceled):
		st = OrderStatusCanceled
	default:
		st = OrderStatusErr
	}

	ttl, _ := money.New(total)

	copied := make([]OrderItem, len(items))
	copy(copied, items)

	return &Order{
		id:        id,
		userID:    userID,
		address:   address,
		items:     copied,
		total:     ttl,
		status:    st,
		createdAt: createdAt,
	}
}

func NewOrder(
	userID int64,
	address string,
	items []OrderItem,
	total money.Money,
) (*Order, error) {
	if userID <= 0 {
		return nil, errors.ErrUserId
	}

	if len(items) == 0 {
		return nil, errors.ErrOrderEmpty
	}

	if strings.TrimSpace(address) == "" {
		return nil, errors.ErrAddressEmpty
	}

	copied := make([]OrderItem, len(items))
	copy(copied, items)

	return &Order{
		userID:    userID,
		address:   strings.TrimSpace(address),
		items:     copied,
		total:     total,
		status:    OrderStatusPendingPayment,
		createdAt: time.Now(),
	}, nil
}

func (o *Order) SwitchStatus(status OrderStatus) {
	o.status = status
}

func (o *Order) Items() []OrderItem {
	items := make([]OrderItem, len(o.items))
	copy(items, o.items)
	return items
}

func NewOrderItem(
	partID, name, brand string,
	quantity int64,
	price money.Money,
	deliveryDay int,
) (*OrderItem, error) {
	if strings.TrimSpace(partID) == "" || len(partID) > 50 {
		return nil, errors.ErrItemPartID
	}

	if strings.TrimSpace(name) == "" || len(name) > 50 {
		return nil, errors.ErrItemName
	}

	if strings.TrimSpace(brand) == "" || len(brand) > 50 {
		return nil, errors.ErrItemBrand
	}

	if quantity <= 0 || quantity > 20 {
		return nil, errors.ErrItemQuantity
	}

	if deliveryDay < 0 || deliveryDay > 1000 {
		return nil, errors.ErrDelivDay
	}

	return &OrderItem{
		partID:      partID,
		name:        name,
		brand:       brand,
		price:       price,
		quantity:    quantity,
		deliveryDay: deliveryDay,
	}, nil
}

func (o *Order) SetID(id int64) {
	o.id = id
}

func (o *Order) UserID() int64 {
	return o.userID
}

func (o *Order) ID() int64 {
	return o.id
}

func (o *Order) Address() string {
	return o.address
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Status() string {
	return string(o.status)
}

func (o *Order) StatusValue() OrderStatus {
	return o.status
}

func (o *Order) Total() money.Money {
	return o.total
}

func (o *OrderItem) PartID() string {
	return o.partID
}

func (o *OrderItem) Brand() string {
	return o.brand
}

func (o *OrderItem) Name() string {
	return o.name
}

func (o *OrderItem) Price() money.Money {
	return o.price
}

func (o *OrderItem) Quantity() int64 {
	return o.quantity
}

func (o *OrderItem) DeliveryDay() int {
	return o.deliveryDay
}
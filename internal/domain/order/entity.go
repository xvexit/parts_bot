package order

import (
	"partsBot/pkg/errors"
	"partsBot/pkg/money"
	"strings"
	"time"
)

type Order struct {
	id        int
	userID    int
	address   string
	items     []OrderItem
	status    OrderStatus
	createdAt time.Time
}

type OrderItem struct {
	partID   string
	name     string
	brand    string
	price    money.Money
	quantity int64
}

func NewOrder(
	userID int,
	address string,
	items []OrderItem,
) (*Order, error) {

	if userID <= 0 {
		return nil, errors.ErrUserId
	}

	if len(items) == 0 {
		return nil, errors.ErrOrderEmpty
	}

	if address == "" {
		return nil, errors.ErrAddressEmpty
	}

	copied := make([]OrderItem, len(items))
	copy(copied, items)

	return &Order{
		userID:    userID,
		address:   address,
		items:     copied,
		status:    OrderStatusNew,
		createdAt: time.Now(),
	}, nil
}

func NewOrderItem(
	partID, name, brand string,
	quantity int64,
	price money.Money,
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

	return &OrderItem{
		partID:   partID,
		name:     name,
		brand:    brand,
		price:    price,
		quantity: quantity,
	}, nil
}

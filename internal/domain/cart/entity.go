package cart

import (
	"partsBot/internal/domain/order"
	"partsBot/pkg/errors"
	"partsBot/pkg/money"
	"strings"
)

type Cart struct {
	userID int
	items  []CartItem
}

type CartItem struct {
	partID   string
	name     string
	brand    string
	price    money.Money
	quantity int64
}

func NewCart(userId int) (*Cart, error) {

	if userId <= 0{
		return nil, errors.ErrUserId
	}

	return &Cart{
		userID: userId,
		items:  make([]CartItem, 0),
	}, nil
}

func (c *Cart) AddItem(item CartItem) {

	for i := range c.items {
		if c.items[i].partID == item.partID{
			if c.items[i].quantity + item.quantity > 20{
				return
			}
			c.items[i].quantity += item.quantity
			return
		}
	}

	c.items = append(c.items, item)
}

func (c *Cart) DeleteItem(partID string) bool {

	for i := range c.items {
		if c.items[i].partID == partID {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Cart) Total() money.Money{
	total, _ := money.New(0)

	for _, item := range c.items {
		total = total.Add(
			item.price.Mul(item.quantity),
		)
	}

	return total
}

func (c *Cart) Clear() {
	c.items = c.items[:0]
}

func NewCartItem(
	partID, name, brand string,
	quantity int64,
	price money.Money,
) (*CartItem, error) {

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

	return &CartItem{
		partID:   partID,
		name:     name,
		brand:    brand,
		price:    price,
		quantity: quantity,
	}, nil
}

func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}

func (c *Cart) Items() []CartItem {
	items := make([]CartItem, len(c.items))
	copy(items, c.items)
	return items
}

func NewOrderFromCart(cart *Cart, address string) (*order.Order, error) {

	if cart.IsEmpty() {
		return nil, errors.ErrCartIsEmpty
	}

	if cart.userID <= 0 {
		return nil, errors.ErrUserId
	}

	items := make([]order.OrderItem, 0, len(cart.items))

	for _, item := range cart.items {
		orderItem, err := order.NewOrderItem(
			item.partID,
			item.name,
			item.brand,
			item.quantity,
			item.price,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, *orderItem)
	}

	return order.NewOrder(cart.userID, address, items)
}

func (c *Cart) ItemsCount() int{
	return len(c.items)
}

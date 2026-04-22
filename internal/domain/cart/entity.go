package cart

import (
	"partsBot/internal/domain/order"
	"partsBot/internal/domain/shared"
	"partsBot/pkg/errors"
	"strings"
)

type Cart struct {
	userID int64
	items  []CartItem
}

type CartItem struct {
	partID      string
	name        string
	brand       string
	price       money.Money
	quantity    int64
	deliveryDay int
	imageURL    string
}

func NewCart(userId int64) (*Cart, error) {

	if userId <= 0 {
		return nil, errors.ErrUserId
	}

	return &Cart{
		userID: userId,
		items:  make([]CartItem, 0),
	}, nil
}

func (c *Cart) AddItem(item CartItem) error {

	for i := range c.items {
		if c.items[i].partID == item.partID {
			if c.items[i].quantity+item.quantity > 20 {
				return errors.ErrItemQuantity
			}
			c.items[i].quantity += item.quantity
			return nil
		}
	}

	c.items = append(c.items, item)
	return nil
}

func (c *Cart) DeleteItem(partID string) error {

	for i := range c.items {
		if c.items[i].partID == partID {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return nil
		}
	}
	return errors.ErrItemNotFound
}

func (c *Cart) Total() money.Money {
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
	partID, name, brand, imageUrl string, //imgurl withiout validation
	quantity int64,
	deliveryDay int,
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

	if deliveryDay < 0 || deliveryDay > 1000 {
		return nil, errors.ErrDelivDay
	}

	return &CartItem{
		partID:      partID,
		name:        name,
		brand:       brand,
		price:       price,
		quantity:    quantity,
		deliveryDay: deliveryDay,
		imageURL:    imageUrl,
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

func (cart *Cart) NewOrderFromCart(address string) (*order.Order, error) {

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
			item.deliveryDay,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, *orderItem)
	}

	return order.NewOrder(cart.userID, address, items)
}

func (c *Cart) ItemsCount() int {
	return len(c.items)
}

func (i CartItem) PartID() string {
	return i.partID
}

func (i CartItem) Name() string {
	return i.name
}

func (i CartItem) Brand() string {
	return i.brand
}

func (i CartItem) Quantity() int64 {
	return i.quantity
}

func (i CartItem) DeliveryDay() int {
	return i.deliveryDay
}

func (i CartItem) ImageURL() string {
	return i.imageURL
}

func (i CartItem) Price() money.Money {
	return i.price
}

func (c Cart) UserID() int64 {
	return c.userID
}

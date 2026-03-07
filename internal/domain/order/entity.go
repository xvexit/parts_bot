package order

import "time"

type Order struct {
	ID        int
	UserID    int
	Adress    string
	Items     []OrderItem
	Status    OrderStatus
	CreatedAt time.Time
}

type OrderItem struct {
	PartID   string
	Name     string
	Brand    string
	Price    int
	Quantity int
}

func NewOrder(
	id, userID int,
	adress string,
	items []OrderItem,
) *Order{
	return &Order{
		ID: id,
		UserID: userID,
		Adress: adress,
		Items: items,
		Status: OrderStatusNew,
		CreatedAt: time.Now(),
	}
}

func NewOrderItem(
	PartID, Name, Brand string,
	Price, Quantity int,
) (*OrderItem, error) {

	if len(PartID) == 0 || len(PartID) > 50 {
		return nil, ErrItemPartID
	}

	if len(Name) == 0 || len(Name) > 50 {
		return nil, ErrItemName
	}

	if len(Brand) == 0 || len(Brand) > 50 {
		return nil, ErrItemBrand
	}

	if Price <= 0 || Price > 1000000 {
		return nil, ErrItemPrice
	}

	if Quantity <= 0 || Quantity > 20 {
		return nil, ErrItemQuantity
	}

	return &OrderItem{
		PartID:   PartID,
		Name:     Name,
		Brand:    Brand,
		Price:    Price,
		Quantity: Quantity,
	}, nil
}
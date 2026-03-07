package cart

type Cart struct {
	UserID int
	Items  []CartItem
}

type CartItem struct {
	PartID   string
	Name     string
	Brand    string
	Price    money.Money
	Quantity int
}

func NewCart(userId int) *Cart {
	return &Cart{
		UserID: userId,
		Items:  make([]CartItem, 0),
	}
}

func (c *Cart) AddItem(item CartItem) {

	for i := range c.Items {
		if c.Items[i].PartID == item.PartID {
			c.Items[i].Quantity += item.Quantity
			return
		}
	}

	c.Items = append(c.Items, item)
}

func (c *Cart) DeleteItem(partID string) {

	for i := range c.Items {
		if c.Items[i].PartID == partID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

func NewCartItem(
	PartID, Name, Brand string,
	Price, Quantity int,
) (*CartItem, error) {

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

	return &CartItem{
		PartID:   PartID,
		Name:     Name,
		Brand:    Brand,
		Price:    Price,
		Quantity: Quantity,
	}, nil
}

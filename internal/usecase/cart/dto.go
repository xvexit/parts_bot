package cart

type CartItemDto struct {
	PartID      string
	Name        string
	Brand       string
	Price       int64
	Quantity    int64
	DeliveryDay int
	ImageURL    string
}

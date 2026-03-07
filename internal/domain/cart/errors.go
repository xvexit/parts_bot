package cart

import "errors"

var(
	ErrItemName = errors.New("product name is too long/short")
	ErrItemBrand = errors.New("brand name is too long/short")
	ErrItemPartID = errors.New("part ID is too long/short")
	ErrItemPrice = errors.New("product price is negative/too high")
	ErrItemQuantity = errors.New("quantity of the product is negative/too large")
)
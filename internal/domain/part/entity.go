package part

import (
	"partsBot/pkg/errors"
	"partsBot/pkg/money"
	"strings"
)

type Part struct {
	ID          int
	Name        string
	Brand       string
	Price       money.Money
	DeliveryDay int
	ImageURL    string
}

func NewPart(
	id, delivDay int,
	name, imgUrl, brand string,
	price money.Money,
) (*Part, error) {

	if id <= 0 {
		return nil, errors.ErrItemPartID
	}

	if strings.TrimSpace(name) == "" || len(name) > 50 {
		return nil, errors.ErrItemName
	}

	if strings.TrimSpace(brand) == "" || len(brand) > 50 {
		return nil, errors.ErrItemBrand
	}

	if delivDay < 0 || delivDay > 20 {
		return nil, errors.ErrItemQuantity
	}

	return &Part{
		ID:          id,
		Name:        name,
		Brand:       brand,
		Price:       price,
		DeliveryDay: delivDay,
		ImageURL:    imgUrl,
	}, nil
}

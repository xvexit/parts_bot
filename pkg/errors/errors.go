package errors

import "errors"

var (
	ErrNameCanNotBeNull   = errors.New("название авто не может быть пустым")
	ErrVinCanNotBeNull    = errors.New("вин код не может быть пустым")
	ErrCarId              = errors.New("car id is too short/long")
	ErrUserId             = errors.New("user id is too short/long")
	ErrItemName           = errors.New("product name is too long/short")
	ErrItemBrand          = errors.New("brand name is too long/short")
	ErrItemPartID         = errors.New("part ID is too long/short")
	ErrItemPrice          = errors.New("product price is negative/too high")
	ErrItemQuantity       = errors.New("quantity of the product is negative/too large")
	ErrUserName           = errors.New("user name is too short/long")
	ErrUserPhone          = errors.New("user phone is too short/long")
	ErrId                 = errors.New("user id is too short/long")
	ErrTgId               = errors.New("user tg id is too short/long")
	ErrAmountCanNotBeNull = errors.New("amount cant be null")
	ErrOrderEmpty         = errors.New("order cant be empty")
	ErrAddressEmpty       = errors.New("address can not be empty")
	ErrCartIsEmpty        = errors.New("cart empty")
	ErrDelivDay           = errors.New("delivery day cant be negative")
)

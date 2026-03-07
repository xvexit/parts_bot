package user

import "errors"

var (
	ErrUserName = errors.New("user name is too short/long")
	ErrUserPhone = errors.New("user phone is too short/long")
	ErrId = errors.New("user id is too short/long")
	ErrTgId = errors.New("user tg id is too short/long")
)

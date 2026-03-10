package user

import "time"

type UserDto struct {
	TelegramID int64
	Name       string
	Phone      string
	CreatedAt  time.Time
}

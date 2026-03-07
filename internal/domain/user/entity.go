package user

import (
	"strings"
	"time"
)

type User struct {
	ID         int
	TelegramID int
	Name       string
	Phone      string
	CreatedAt  time.Time
}

func NewUser(
	id, tgid int,
	name, phone string,
) (*User, error) {

	if id <= 0 || id > 50 {
		return nil, ErrId
	}

	if id <= 0 || id > 50 {
		return nil, ErrTgId
	}

	if strings.TrimSpace(name) == "" || len(name) > 50 {
		return nil, ErrUserName
	}

	if strings.TrimSpace(phone) == "" || len(phone) > 12 {
		return nil, ErrUserPhone
	}

	return &User{
		ID:         id,
		TelegramID: tgid,
		Name:       name,
		Phone:      phone,
		CreatedAt:  time.Now(),
	}, nil
}

func (u *User) ChangePhone(newPhone string) error {
	
	if strings.TrimSpace(newPhone) == "" || len(newPhone) > 12 {
		return ErrUserPhone
	}

	u.Phone = newPhone

	return nil
}

func (u *User) ChangeName(newName string) error {
	
	if strings.TrimSpace(newName) == "" || len(newName) > 12 {
		return ErrUserName
	}

	u.Name = newName

	return nil
}

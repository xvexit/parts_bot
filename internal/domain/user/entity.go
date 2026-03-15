package user

import (
	"partsBot/pkg/errors"
	"strings"
	"time"
)

type User struct {
	id         int64
	telegramID int64
	name       string
	phone      string
	createdAt  time.Time
}

func NewUser(
	tgid int64,
	name, phone string,
) (*User, error) {

	if tgid <= 0 {
		return nil, errors.ErrTgId
	}

	if strings.TrimSpace(name) == "" || len(name) > 50 {
		return nil, errors.ErrUserName
	}

	if strings.TrimSpace(phone) == "" || len(phone) > 20 {
		return nil, errors.ErrUserPhone
	}

	return &User{
		telegramID: tgid,
		name:       name,
		phone:      phone,
		createdAt:  time.Now(),
	}, nil
}

func (u *User) ChangePhone(newPhone string) error {

	if strings.TrimSpace(newPhone) == "" || len(newPhone) > 20 {
		return errors.ErrUserPhone
	}

	u.phone = newPhone

	return nil
}

func (u *User) ChangeName(newName string) error {

	if strings.TrimSpace(newName) == "" || len(newName) > 50 {
		return errors.ErrUserName
	}

	u.name = newName

	return nil
}

func RestoreUser(
	id int64,
	telegramID int64,
	name, phone string,
	createdAt time.Time,
) *User {

	return &User{
		id:         id,
		telegramID: telegramID,
		name:       name,
		phone:      phone,
		createdAt:  createdAt,
	}
}

func (u *User) SetID(id int64) {
	u.id = id
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) TelegramID() int64 {
	return u.telegramID
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Phone() string {
	return u.phone
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

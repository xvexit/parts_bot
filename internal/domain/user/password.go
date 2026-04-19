package user

import (
	"partsBot/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash string
}

const minPasswordLength = 6

// New создаёт Password из сырого пароля (регистрация)
func NewPassword(raw string) (Password, error) {
	if len(raw) < minPasswordLength {
		return Password{}, errors.ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(hash)}, nil
}

// FromHash используется при загрузке из БД
func PasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Compare проверяет пароль (логин)
func (p Password) Compare(raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(raw))
	return err == nil
}

// Hash возвращает хэш для сохранения в БД
func (p Password) Hash() string {
	return p.hash
}

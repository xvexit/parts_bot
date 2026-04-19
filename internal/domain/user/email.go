package user

import (
	"regexp"
	"strings"

	"partsBot/pkg/errors"
)

type Email struct {
	value string
}

// упрощённая, но практичная regex
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// NewEmail создаёт email из сырой строки (регистрация / обновление)
func NewEmail(raw string) (Email, error) {
	normalized := normalizeEmail(raw)

	if normalized == "" {
		return Email{}, errors.ErrInvalidEmail
	}

	if !emailRegex.MatchString(normalized) {
		return Email{}, errors.ErrInvalidEmail
	}

	return Email{value: normalized}, nil
}

// EmailFromDB используется при чтении из БД
func EmailFromDB(value string) (Email, error) {
	// можно валидировать, можно нет — я рекомендую валидировать
	if value == "" {
		return Email{}, errors.ErrInvalidEmail
	}

	if !emailRegex.MatchString(value) {
		return Email{}, errors.ErrInvalidEmail
	}

	return Email{value: value}, nil
}

// String возвращает email (для БД / ответа)
func (e Email) String() string {
	return e.value
}

// Equal сравнение email
func (e Email) Equal(other Email) bool {
	return e.value == other.value
}

// normalizeEmail приводит к каноническому виду
func normalizeEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

package money

import (
	"partsBot/pkg/errors"
)

type Money struct {
	amount int64
}

func New(amount int64) (Money, error) {

    if amount < 0{
        return Money{}, errors.ErrAmountCanNotBeNull
    }

	return Money{amount: amount}, nil
}

func (m Money) Add(other Money) Money {
    return Money{amount: m.amount + other.amount}
}

func (m Money) Mul(qty int64) Money {
    return Money{amount: m.amount * qty}
}

func (m Money) Amount() int64 {
    return m.amount
}
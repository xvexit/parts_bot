package car

import (
	"partsBot/pkg/errors"
	"strings"
)

type Car struct {
	ID     int
	Name   string
	VIN    string //vin
	UserID int
}

func NewCar(id, userID int, name, vin string) (*Car, error) {

	if id <= 0 {
		return nil, errors.ErrId
	}

	if userID <= 0 {
		return nil, errors.ErrUserId
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.ErrNameCanNotBeNull
	}

	if strings.TrimSpace(vin) == "" {
		return nil, errors.ErrVinCanNotBeNull
	}

	return &Car{
		ID:     id,
		Name:   name,
		VIN:    vin,
		UserID: userID,
	}, nil
}

func (c *Car) ChangeName(newName string) error {

	if strings.TrimSpace(newName) == "" {
		return errors.ErrNameCanNotBeNull
	}

	c.Name = newName

	return nil
}

func (c *Car) ChangeVin(newVin string) error {

	if strings.TrimSpace(newVin) == "" {
		return errors.ErrVinCanNotBeNull
	}

	c.VIN = newVin

	return nil
}

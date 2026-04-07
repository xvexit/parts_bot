package car

import (
	"partsBot/pkg/errors"
	"strings"
)

type Car struct {
	id     int64
	name   string
	vin    string //vin
	userID int64
}

func NewCar(userID int64, name, vin string) (*Car, error) {

	if userID <= 0 {
		return nil, errors.ErrUserId
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.ErrNameCanNotBeNull
	}

	if len(strings.TrimSpace(vin)) != 17 {
		return nil, errors.ErrVinMustBe17
	}

	return &Car{
		name:   name,
		vin:    vin,
		userID: userID,
	}, nil
}

func RestoreCar(id, userID int64, name, vin string) *Car {
	return &Car{
		id:     id,
		userID: userID,
		name:   name,
		vin:    vin,
	}
}

func (c *Car) SetId(id int64) {
	c.id = id
}

func (c *Car) ChangeName(newName string) error {

	if strings.TrimSpace(newName) == "" {
		return errors.ErrNameCanNotBeNull
	}

	c.name = newName

	return nil
}

func (c *Car) ChangeVin(newVin string) error {

	if len(strings.TrimSpace(newVin)) != 17 {
		return errors.ErrVinMustBe17
	}

	c.vin = newVin

	return nil
}

func (c *Car) UserId() int64 {
	return c.userID
}

func (c *Car) Vin() string {
	return c.vin
}

func (c *Car) Name() string {
	return c.name
}

func (c *Car) ID() int64 {
	return c.id
}

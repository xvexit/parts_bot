package car

import "strings"

type Car struct {
	Id     int
	Name   string
	Vin    string //vin
	UserID int
}

func NewCar(id, userID int, name, vin string) (*Car, error) {

	if strings.TrimSpace(name) == "" {
		return nil, ErrNameCanNotBeNull
	}

	if strings.TrimSpace(vin) == "" {
		return nil, ErrVinCanNotBeNull
	}

	return &Car{
		Id:     id,
		Name:   name,
		Vin:    vin,
		UserID: userID,
	}, nil
}

func (c *Car) ChangeName(newName string) error {

	if strings.TrimSpace(newName) == "" {
		return ErrNameCanNotBeNull
	}

	c.Name = newName

	return nil
}

func (c *Car) ChangeVin(newVin string) error {

	if strings.TrimSpace(newVin) == "" {
		return ErrVinCanNotBeNull
	}

	c.Vin = newVin

	return nil
}

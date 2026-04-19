package car

import (
	"context"
	"testing"

	domain "partsBot/internal/domain/car"
)

type stubCarRepo struct {
	createFn func(ctx context.Context, c *domain.Car) error
	updateFn func(ctx context.Context, c *domain.Car) error
	getFn    func(ctx context.Context, id int64) (*domain.Car, error)
	deleteFn func(ctx context.Context, id int64) error
	listFn   func(ctx context.Context, userID int64) ([]domain.Car, error)
}

func (s *stubCarRepo) Create(ctx context.Context, c *domain.Car) error {
	return s.createFn(ctx, c)
}

func (s *stubCarRepo) Update(ctx context.Context, c *domain.Car) error {
	return s.updateFn(ctx, c)
}

func (s *stubCarRepo) GetByID(ctx context.Context, id int64) (*domain.Car, error) {
	return s.getFn(ctx, id)
}

func (s *stubCarRepo) Delete(ctx context.Context, id int64) error {
	return s.deleteFn(ctx, id)
}

func (s *stubCarRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Car, error) {
	return s.listFn(ctx, userID)
}

func TestAddCar(t *testing.T) {

	repo := &stubCarRepo{
		createFn: func(ctx context.Context, c *domain.Car) error {
			c.SetId(1)
			return nil
		},
	}

	service := Service{repo: repo}

	dto := CarInput{
		UserID: 1,
		Name:   "BMW",
		VIN:    "123",
	}

	car, err := service.Add(context.Background(), dto)

	if err != nil {
		t.Fatal(err)
	}

	if car.ID() != 1 {
		t.Fatal("id not set")
	}
}

func TestChangeCarName(t *testing.T) {

	car := domain.RestoreCar(1, 1, "BMW", "VIN")

	repo := &stubCarRepo{
		getFn: func(ctx context.Context, id int64) (*domain.Car, error) {
			return car, nil
		},
		updateFn: func(ctx context.Context, c *domain.Car) error {
			return nil
		},
	}

	service := Service{repo: repo}

	res, err := service.ChangeCar(context.Background(), 1, "Audi")

	if err != nil {
		t.Fatal(err)
	}

	if res.Name() != "Audi" {
		t.Fatal("name not updated")
	}
}

func TestDeleteCar(t *testing.T) {

	repo := &stubCarRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	service := Service{repo: repo}

	err := service.Delete(context.Background(), 1)

	if err != nil {
		t.Fatal(err)
	}
}

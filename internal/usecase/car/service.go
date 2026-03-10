package car

import (
	"context"
	"partsBot/internal/domain/car"
)

type Service struct {
	repo car.Repository
}

func (s *Service) Add(ctx context.Context, carDto CarDto) (car.Car, error) { // защитить от создания множества машин на 1 акк

	ccar, err := car.NewCar(carDto.UserID, carDto.Name, carDto.VIN)
	if err != nil {
		return car.Car{}, err
	}

	if err := s.repo.Save(ctx, ccar); err != nil {
		return car.Car{}, err
	}

	return *ccar, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id int64) (car.Car, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByUser(ctx context.Context, userId int64) ([]car.Car, error) {
	return s.repo.ListByUser(ctx, userId)
}

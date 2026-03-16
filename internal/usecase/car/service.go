package car

import (
	"context"
	"partsBot/internal/domain/car"
)

type Service struct {
	repo car.Repository
}

func NewService(repo car.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, carDto CarDto) (*car.Car, error) { // защитить от создания множества машин на 1 акк

	ccar, err := car.NewCar(carDto.UserID, carDto.Name, carDto.VIN)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, ccar); err != nil {
		return nil, err
	}

	return ccar, nil
}

func (s *Service) ChangeCar(ctx context.Context, id int64, newName string) (*car.Car, error) {
	car, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := car.ChangeName(newName); err != nil{
		return nil, err
	}

	if err := s.repo.Update(ctx, car); err != nil{
		return nil, err
	}

	return car, nil
}

func (s *Service) ChangeVin(ctx context.Context, id int64, newVin string) (*car.Car, error) {
	car, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := car.ChangeVin(newVin); err != nil{
		return nil, err
	}

	if err := s.repo.Update(ctx, car); err != nil{
		return nil, err
	}

	return car, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*car.Car, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByUser(ctx context.Context, userId int64) ([]car.Car, error) {
	return s.repo.ListByUser(ctx, userId)
}

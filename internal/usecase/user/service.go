package user

import (
	"context"
	"partsBot/internal/domain/user"
)

type Service struct {
	repo user.Repository
}

func (s *Service) Register(
    ctx context.Context,
    dto UserDto,
) (*user.User, error){
	userr, err := user.NewUser(dto.TelegramID, dto.Name, dto.Phone)
	if err != nil{
		return nil, err
	}

	if err := s.repo.Create(ctx, userr); err != nil{
		return nil, err
	}

	return userr, nil
}

func (s *Service) ChangeName(ctx context.Context, newName string, userID int64) (*user.User, error) {
	userr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := userr.ChangeName(newName); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, userr); err != nil{
		return nil, err
	}

	return userr, err
}

func (s *Service) ChangePhone(ctx context.Context, newPhone string, userID int64) (*user.User, error) {
	userr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := userr.ChangePhone(newPhone); err != nil {
		return nil, err
	}

	return userr, s.repo.Update(ctx, userr)
}

func (s *Service) Delete(ctx context.Context, userID int64) error {
	return s.repo.Delete(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, userID int64) (*user.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *Service) GetByTgID(ctx context.Context, userTgID int64) (*user.User, error) {
	return s.repo.GetByTgID(ctx, userTgID)
}
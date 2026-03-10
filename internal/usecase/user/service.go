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

	if err := s.repo.Save(ctx, userr); err != nil{
		return nil, err
	}

	return userr, nil
}

func (s *Service) ChangeName(ctx context.Context, newName string, userID int) (user.User, error) {
	userr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return user.User{}, err
	}

	if err := userr.ChangeName(newName); err != nil {
		return user.User{}, err
	}

	return userr, s.repo.Save(ctx, &userr)
}

func (s *Service) ChangePhone(ctx context.Context, newPhone string, userID int) (user.User, error) {
	userr, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return user.User{}, err
	}

	if err := userr.ChangePhone(newPhone); err != nil {
		return user.User{}, err
	}

	return userr, s.repo.Save(ctx, &userr)
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	return s.repo.Delete(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, userID int) (user.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *Service) GetByTgID(ctx context.Context, userTgID int) (user.User, error) {
	return s.repo.GetByTgID(ctx, userTgID)
}
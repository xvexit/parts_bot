package auth

import (
	"context"
	"partsBot/internal/domain/user"
	"partsBot/pkg/errors"
)

type Service struct {
	tm       TokenManager
	userRepo user.Repository
}

func NewService(tm TokenManager, repo user.Repository) *Service {
	return &Service{
		userRepo: repo,
		tm:       tm,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return "", "", err
	}

	u, err := s.userRepo.GetByEmail(ctx, emailVO.String())
	if err != nil {
		return "", "", err
	}

	if !u.Pass().Compare(password) {
		return "", "", errors.ErrInvalidCredentials
	}

	access, err := s.tm.GenerateAccessToken(u.ID())
	if err != nil {
		return "", "", err
	}

	refresh, err := s.tm.GenerateRefreshToken(u.ID())
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.tm.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	return s.tm.GenerateAccessToken(userID)
}

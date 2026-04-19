	package user

	import (
		"context"
		"partsBot/internal/domain/user"
	)

	type Service struct {
		repo user.Repository
	}

	func NewService(repo user.Repository) *Service {
		return &Service{
			repo: repo,
		}
	}

	func (s *Service) Register(
		ctx context.Context,
		dto UserInput,
	) (*user.User, error) {

		pass, err := user.NewPassword(dto.Password)
		if err != nil {
			return nil, err
		}

		email, err := user.NewEmail(dto.Email) //баг пустой емаил не создастся (может сделать чтобы он возвращал укаpатель?)
		if err != nil {
			return nil, err
		}

		userr, err := user.NewUser(dto.Name, dto.Phone, pass, email)
		if err != nil {
			return nil, err
		}

		return s.repo.Create(ctx, userr);
	}

	func (s *Service) ChangeName(ctx context.Context, newName string, userID int64) (*user.User, error) {
		userr, err := s.repo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if err := userr.ChangeName(newName); err != nil {
			return nil, err
		}

		return s.repo.Update(ctx, userr)
	}

	func (s *Service) ChangePhone(ctx context.Context, newPhone string, userID int64) (*user.User, error) {
		userr, err := s.repo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if err := userr.ChangePhone(newPhone); err != nil {
			return nil, err
		}

		return s.repo.Update(ctx, userr)
	}

	func (s *Service) Delete(ctx context.Context, userID int64) error {
		return s.repo.Delete(ctx, userID)
	}

	func (s *Service) GetByID(ctx context.Context, userID int64) (*user.User, error) {
		return s.repo.GetByID(ctx, userID)
	}

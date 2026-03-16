package order

import (
	"context"
	"partsBot/internal/domain/order"
)

type Service struct {
	repo order.Repository
}

func NewService(repo order.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByUserID(ctx context.Context, userID int64) ([]order.Order, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) OrderItems(ctx context.Context, id int64) ([]order.OrderItem, error) {

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order.Items(), nil
}

func (s *Service) SwitchStatus(ctx context.Context, id int64, status order.OrderStatus) error {

	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	ord.SwitchStatus(status)

	return s.repo.UpdateStatus(ctx, id, status)
}

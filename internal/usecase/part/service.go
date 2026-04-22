package part

import (
	"context"
	"strings"
)

type CatalogItem struct {
	PartID      string
	Name        string
	Brand       string
	Price       int64
	DeliveryDay int
}

type CatalogRepository interface {
	Search(ctx context.Context, query string, limit int) ([]CatalogItem, error)
}

type Service struct {
	repo CatalogRepository
}

func NewService(repo CatalogRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]CatalogItem, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return []CatalogItem{}, nil
	}
	if limit <= 0 {
		limit = 30
	}
	return s.repo.Search(ctx, q, limit)
}

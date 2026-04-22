package repository

import (
	"context"
	"fmt"
	"partsBot/internal/infrastructure/db"
	partuc "partsBot/internal/usecase/part"
)

type PostgresPartCatalogRepo struct {
	db *db.DB
}

func NewPostgresPartCatalogRepo(db *db.DB) *PostgresPartCatalogRepo {
	return &PostgresPartCatalogRepo{db: db}
}

func (r *PostgresPartCatalogRepo) Search(ctx context.Context, query string, limit int) ([]partuc.CatalogItem, error) {
	exec := r.db.Executor(ctx)
	rows, err := exec.Query(
		ctx,
		`SELECT part_id, name, brand, price, delivery_day
		 FROM catalog_parts
		 WHERE part_id ILIKE $1 OR name ILIKE $1 OR brand ILIKE $1
		 ORDER BY
		     CASE WHEN part_id ILIKE $2 THEN 0 ELSE 1 END,
		     name ASC
		 LIMIT $3`,
		"%"+query+"%",
		query+"%",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("search catalog query failed: %w", err)
	}
	defer rows.Close()

	result := make([]partuc.CatalogItem, 0, limit)
	for rows.Next() {
		var item partuc.CatalogItem
		if err := rows.Scan(&item.PartID, &item.Name, &item.Brand, &item.Price, &item.DeliveryDay); err != nil {
			return nil, fmt.Errorf("search catalog scan failed: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search catalog rows failed: %w", err)
	}
	return result, nil
}

package repository

import (
	"context"

	"partsBot/internal/domain/order"
	money "partsBot/internal/domain/shared"
	"partsBot/internal/infrastructure/db"
)

type PostgresOrderRepo struct {
	db *db.DB
}

func NewPostgresOrderRepo(db *db.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{
		db: db,
	}
}

func (r *PostgresOrderRepo) Create(ctx context.Context, ord *order.Order) error {
	exec := r.db.Executor(ctx)

	query := `
	INSERT INTO orders(user_id, address, status, total, created_at)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		ord.UserID(),
		ord.Address(),
		ord.Status(),
		ord.Total().Amount(),
		ord.CreatedAt(),
	).Scan(&id)
	if err != nil {
		return err
	}

	ord.SetID(id)

	for _, item := range ord.Items() {
		_, err := exec.Exec(
			ctx,
			`INSERT INTO order_items(order_id, part_id, name, brand, price, quantity, delivery_day)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			id,
			item.PartID(),
			item.Name(),
			item.Brand(),
			item.Price().Amount(),
			item.Quantity(),
			item.DeliveryDay(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresOrderRepo) GetByID(ctx context.Context, orderID int64) (*order.Order, error) {
	exec := r.db.Executor(ctx)

	query := `
	SELECT id, user_id, address, status, total, created_at
	FROM orders
	WHERE id=$1
	`

	model := OrderModel{}

	err := exec.QueryRow(ctx, query, orderID).Scan(
		&model.ID,
		&model.UserID,
		&model.Address,
		&model.Status,
		&model.Total,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	items, err := r.getItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return order.RestoreOrder(
		model.ID,
		model.UserID,
		model.Total,
		items,
		model.Status,
		model.Address,
		model.CreatedAt,
	), nil
}

func (r *PostgresOrderRepo) getItems(ctx context.Context, orderID int64) ([]order.OrderItem, error) {
	exec := r.db.Executor(ctx)

	rows, err := exec.Query(
		ctx,
		`SELECT part_id, name, brand, price, quantity, delivery_day
		 FROM order_items
		 WHERE order_id=$1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemSl := make([]order.OrderItem, 0)

	for rows.Next() {
		model := OrderItemModel{}

		err := rows.Scan(
			&model.PartID,
			&model.Name,
			&model.Brand,
			&model.Price,
			&model.Quantity,
			&model.DeliveryDay,
		)
		if err != nil {
			return nil, err
		}

		price, err := money.New(model.Price)
		if err != nil {
			return nil, err
		}

		item, err := order.NewOrderItem(
			model.PartID,
			model.Name,
			model.Brand,
			model.Quantity,
			price,
			model.DeliveryDay,
		)
		if err != nil {
			return nil, err
		}

		itemSl = append(itemSl, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return itemSl, nil
}

func (r *PostgresOrderRepo) ListByUserID(ctx context.Context, userID int64) ([]order.Order, error) {
	exec := r.db.Executor(ctx)

	query := `
	SELECT id, user_id, address, total, status, created_at
	FROM orders
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, err := exec.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderSl []order.Order

	for rows.Next() {
		model := OrderModel{}

		err := rows.Scan(
			&model.ID,
			&model.UserID,
			&model.Address,
			&model.Total,
			&model.Status,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items, err := r.getItems(ctx, model.ID)
		if err != nil {
			return nil, err
		}

		ord := order.RestoreOrder(
			model.ID,
			model.UserID,
			model.Total,
			items,
			model.Status,
			model.Address,
			model.CreatedAt,
		)

		orderSl = append(orderSl, *ord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderSl, nil
}

func (r *PostgresOrderRepo) UpdateStatus(
	ctx context.Context,
	id int64,
	status order.OrderStatus,
) error {
	exec := r.db.Executor(ctx)

	_, err := exec.Exec(
		ctx,
		`UPDATE orders
		 SET status=$1
		 WHERE id=$2`,
		status,
		id,
	)

	return err
}
package repository

import (
	"context"
	"partsBot/internal/domain/order"
	"partsBot/internal/infrastructure/db"
	"partsBot/pkg/money"
)

type PostgresOrderRepo struct {
	db *db.DB
}

func NewPostgresOrderRepo(db *db.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{
		db: db,
	}
}

func (r *PostgresOrderRepo) Create(ctx context.Context, order *order.Order) error {
	exec := r.db.Executor(ctx)

	query := `
	INSERT INTO orders(user_id, address, status, created_at)
	VALUES ($1,$2,$3,$4)
	RETURNING id
	`
	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		order.UserID(),
		order.Address(),
		order.Status(),
		order.CreatedAt(),
	).Scan(&id)

	if err != nil {
		return err
	}

	order.SetID(id)

	for _, item := range order.Items() {
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
	SELECT id, user_id, address, status, created_at
	FROM orders
	WHERE id=$1
	`

	dto := OrderDTO{}

	err := exec.QueryRow(
		ctx,
		query,
		orderID,
	).Scan(
		&dto.ID,
		&dto.UserID,
		&dto.Address,
		&dto.Status,
		&dto.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	items, err := r.getItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order := order.RestoreOrder(
		dto.ID,
		dto.UserID,
		dto.Address,
		items,
		dto.Status,
		dto.CreatedAt,
	)

	return order, nil
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

	itemSl := []order.OrderItem{}

	for rows.Next() {
		dto := OrderItemDTO{}

		err := rows.Scan(
			&dto.PartID,
			&dto.Name,
			&dto.Brand,
			&dto.Price,
			&dto.Quantity,
			&dto.DeliveryDay,
		)

		if err != nil {
			return nil, err
		}

		price, err := money.New(dto.Price)
		if err != nil {
			return nil, err
		}

		item, err := order.NewOrderItem(
			dto.PartID,
			dto.Name,
			dto.Brand,
			dto.Quantity,
			price,
			dto.DeliveryDay,
		)
		if err != nil {
			return nil, err
		}

		itemSl = append(itemSl, *item)
	}

	return itemSl, nil
}

func (r *PostgresOrderRepo) ListByUserID(ctx context.Context, userID int64) ([]order.Order, error) {
	exec := r.db.Executor(ctx)

	query := `
	SELECT id, user_id, address, status, created_at
	FROM orders
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, err := exec.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderSl []order.Order

	for rows.Next() {
		dto := OrderDTO{}
		err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Address,
			&dto.Status,
			&dto.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items, err := r.getItems(ctx, dto.ID)
		if err != nil {
			return nil, err
		}

		ord := order.RestoreOrder(
			dto.ID,
			dto.UserID,
			dto.Address,
			items,
			dto.Status,
			dto.CreatedAt,
		)

		orderSl = append(orderSl, *ord)
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

package repository

import (
	"context"
	"partsBot/internal/domain/cart"
	"partsBot/internal/infrastructure/db"
	"partsBot/pkg/money"
)

type PostgresCartRepo struct {
	db *db.DB
}

func NewPostgresCartRepo(db *db.DB) *PostgresCartRepo {
	return &PostgresCartRepo{
		db: db,
	}
}

func (r *PostgresCartRepo) Save(ctx context.Context, cart *cart.Cart) error {
	exec := r.db.Executor(ctx)

	_, err := exec.Exec(
		ctx,
		`DELETE FROM cart_items WHERE user_id = $1`,
		cart.UserID(),
	)
	if err != nil {
		return err
	}

	for _, item := range cart.Items() {
		_, err := exec.Exec(
			ctx,
			`INSERT INTO cart_items
			(user_id, part_id, name, brand, price, quantity, delivery_day, image_url)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			cart.UserID(),
			item.PartID(),
			item.Name(),
			item.Brand(),
			item.Price().Amount(),
			item.Quantity(),
			item.DeliveryDay(),
			item.ImageURL(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresCartRepo) GetByUserID(ctx context.Context, userID int64) (*cart.Cart, error) {
	exec := r.db.Executor(ctx)

	rows, err := exec.Query(
		ctx,
		`SELECT part_id, name, brand, price, quantity, delivery_day, image_url
		FROM cart_items
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cartAgg, err := cart.NewCart(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dto := CartItemDto{}

		err := rows.Scan(
			&dto.PartID,
			&dto.Name,
			&dto.Brand,
			&dto.Price,
			&dto.Quantity,
			&dto.DeliveryDay,
			&dto.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		moneyPrice, err := money.New(dto.Price)
		if err != nil {
			return nil, err
		}

		item, err := cart.NewCartItem(
			dto.PartID,
			dto.Name,
			dto.Brand,
			dto.ImageURL,
			dto.Quantity,
			dto.DeliveryDay,
			moneyPrice,
		)
		if err != nil {
			return nil, err
		}

		err = cartAgg.AddItem(*item)
		if err != nil {
			return nil, err
		}
	}
	return cartAgg, nil
}

func (r *PostgresCartRepo) Delete(ctx context.Context, userID int64) error {
	exec := r.db.Executor(ctx)

	_, err := exec.Exec(
		ctx,
		`DELETE FROM cart_items WHERE user_id=$1`,
		userID,
	)

	return err
}

func (r *PostgresCartRepo) GetByUserIDForUpdate(
	ctx context.Context,
	userID int64,
) (*cart.Cart, error) {

	exec := r.db.Executor(ctx)

	rows, err := exec.Query(
		ctx,
		`SELECT part_id, name, brand, price, quantity, delivery_day, image_url
		FROM cart_items
		WHERE user_id = $1
		FOR UPDATE`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cartAgg, err := cart.NewCart(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		dto := CartItemDto{}

		err := rows.Scan(
			&dto.PartID,
			&dto.Name,
			&dto.Brand,
			&dto.Price,
			&dto.Quantity,
			&dto.DeliveryDay,
			&dto.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		price, err := money.New(dto.Price)
		if err != nil {
			return nil, err
		}

		item, err := cart.NewCartItem(
			dto.PartID,
			dto.Name,
			dto.Brand,
			dto.ImageURL,
			dto.Quantity,
			dto.DeliveryDay,
			price,
		)
		if err != nil {
			return nil, err
		}

		cartAgg.AddItem(*item)
	}

	return cartAgg, nil
}
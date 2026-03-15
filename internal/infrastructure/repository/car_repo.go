package repository

import (
	"context"

	"partsBot/internal/domain/car"
	"partsBot/internal/infrastructure/db"
)

type PostgresCarRepo struct {
	db *db.DB
}

func NewPostgresCarRepo(db *db.DB) *PostgresCarRepo {
	return &PostgresCarRepo{
		db: db,
	}
}

func (r *PostgresCarRepo) Create(ctx context.Context, c *car.Car) error {

	exec := r.db.Executor(ctx)

	query := `
	INSERT INTO cars(user_id, name, vin)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		c.UserId(),
		c.Name(),
		c.Vin(),
	).Scan(&id)

	if err != nil {
		return err
	}

	c.SetId(id)

	return nil
}

func (r *PostgresCarRepo) Update(ctx context.Context, c *car.Car) error {

	exec := r.db.Executor(ctx)

	query := `
	UPDATE cars SET name=$1, vin=$2
	WHERE id=$3
	`

	_, err := exec.Exec(
		ctx,
		query,
		c.Name(),
		c.Vin(),
		c.ID(),
	)

	return err
}

func (r *PostgresCarRepo) GetByID(ctx context.Context, id int64) (*car.Car, error) {
	exec := r.db.Executor(ctx)

	query := `
	SELECT id, user_id, name, vin
	FROM cars
	WHERE id = $1
	`

	dto := CarDTO{}

	err := exec.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&dto.ID,
		&dto.UserID,
		&dto.Name,
		&dto.VIN,
	)

	if err != nil {
		return nil, err
	}

	car := car.RestoreCar(
		dto.ID,
		dto.UserID,
		dto.Name,
		dto.VIN,
	)
	
	return car, nil
}

func (r *PostgresCarRepo) ListByUser(ctx context.Context, userID int64) ([]car.Car, error) {

	exec := r.db.Executor(ctx)

	query := `
	SELECT id, user_id, name, vin 
	FROM cars
	WHERE user_id=$1
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

	var carSl []car.Car

	for rows.Next() {
		dto := CarDTO{}

		err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Name,
			&dto.VIN,
		)

		if err != nil {
			return nil, err
		}

		car := car.RestoreCar(
			dto.ID,
			dto.UserID,
			dto.Name,
			dto.VIN,
		)

		carSl = append(carSl, *car)
	}

	if err := rows.Err(); err != nil{
		return nil, err
	}

	return carSl, nil
}

func (r *PostgresCarRepo) Delete(ctx context.Context, id int64) error {
	exec := r.db.Executor(ctx)

	query := `
	DELETE FROM cars
	WHERE id=$1
	`

	_, err := exec.Exec(
		ctx,
		query,
		id,
	)

	return err

}

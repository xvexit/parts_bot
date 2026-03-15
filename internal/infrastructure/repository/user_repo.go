package repository

import (
	"context"
	"partsBot/internal/domain/user"
	"partsBot/internal/infrastructure/db"
)

type PostgresUserRepo struct {
	db *db.DB
}

func NewPostgresUserRepo(db *db.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		db: db,
	}
}

func (r *PostgresUserRepo) Create(ctx context.Context, us *user.User) error {

	exec := r.db.Executor(ctx)

	query := `
	INSERT INTO users(telegram_id, name, phone, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING id
	`

	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		us.TelegramID(),
		us.Name(),
		us.Phone(),
		us.CreatedAt(),
	).Scan(&id)

	if err != nil{
		return err
	}

	us.SetID(id)

	return nil
}

func (r *PostgresUserRepo) Update(ctx context.Context, us *user.User) error {

	exec := r.db.Executor(ctx)

	query := `
	UPDATE users SET name=$1, phone=$2
	WHERE id=$3
	`

	_, err := exec.Exec(
		ctx,
		query,
		us.Name(),
		us.Phone(),
		us.ID(),
	)

	return err
}

func (r *PostgresUserRepo) Delete(ctx context.Context, userID int64) error{
	exec := r.db.Executor(ctx)

	query := `
	DELETE FROM users
	WHERE id=$1
	`

	_, err := exec.Exec(
		ctx,
		query,
		userID,
	)
	return err
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, userID int64) (*user.User, error) {

	exec := r.db.Executor(ctx)

	query := `
	SELECT id, telegram_id, name, phone, created_at
	FROM users
	WHERE id=$1
	`

	dto := UserDTO{}

	err := exec.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&dto.ID,
		&dto.TelegramID,
		&dto.Name,
		&dto.Phone,
		&dto.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user.RestoreUser(
		dto.ID,
		dto.TelegramID,
		dto.Name,
		dto.Phone,
		dto.CreatedAt,
	), nil	
}

func (r *PostgresUserRepo) GetByTgID(ctx context.Context, userTgID int64) (*user.User, error){

	exec := r.db.Executor(ctx)

	query := `
	SELECT id, telegram_id, name, phone, created_at
	FROM users
	WHERE telegram_id=$1
	`

	dto := UserDTO{}

	err := exec.QueryRow(
		ctx,
		query,
		userTgID,
	).Scan(
		&dto.ID,
		&dto.TelegramID,
		&dto.Name,
		&dto.Phone,
		&dto.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user.RestoreUser(
		dto.ID,
		dto.TelegramID,
		dto.Name,
		dto.Phone,
		dto.CreatedAt,
	), nil	
}

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

func (r *PostgresUserRepo) Create(ctx context.Context, us *user.User) (*user.User, error) {
	exec := r.db.Executor(ctx)

	query := `
	INSERT INTO users(name, email, phone, password_hash, created_at)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id
	`

	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		us.Name(),
		us.Email(),
		us.Phone(),
		us.Pass().Hash(),
		us.CreatedAt(),
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	us.SetID(id)
	return us, nil
}

func (r *PostgresUserRepo) Update(ctx context.Context, us *user.User) (*user.User, error) {
	exec := r.db.Executor(ctx)

	query := `
	UPDATE users 
	SET name=$1, phone=$2, email=$3, password_hash=$4
	WHERE id=$5
	`
	_, err := exec.Exec(
		ctx,
		query,
		us.Name(),
		us.Phone(),
		us.Email(),
		us.Pass().Hash(),
		us.ID(),
	)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (r *PostgresUserRepo) Delete(ctx context.Context, userID int64) error {
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
	SELECT id, name, email, phone, password_hash, created_at
	FROM users
	WHERE id=$1
	`

	row := exec.QueryRow(ctx, query, userID)

	return scanUser(row)
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	exec := r.db.Executor(ctx)

	query := `
	SELECT id, name, email, phone, password_hash, created_at
	FROM users
	WHERE email = $1
	`

	row := exec.QueryRow(ctx, query, email)

	u, err := scanUser(row)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// func (r *PostgresUserRepo) GetByTgID(ctx context.Context, userTgID int64) (*user.User, error) {

// 	exec := r.db.Executor(ctx)

// 	query := `
// 	SELECT id, telegram_id, name, phone, created_at
// 	FROM users
// 	WHERE telegram_id=$1
// 	`

// 	dto := UserDTO{}

// 	err := exec.QueryRow(
// 		ctx,
// 		query,
// 		userTgID,
// 	).Scan(
// 		&dto.ID,
// 		&dto.TelegramID,
// 		&dto.Name,
// 		&dto.Phone,
// 		&dto.CreatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return user.RestoreUser(
// 		dto.ID,
// 		dto.TelegramID,
// 		dto.Name,
// 		dto.Phone,
// 		dto.CreatedAt,
// 	), nil
// }

func scanUser(row scanner) (*user.User, error) {
	var m UserModel

	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Email,
		&m.Phone,
		&m.Password,
		&m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	var emailVO user.Email
	if m.Email != nil {
		e, err := user.EmailFromDB(*m.Email)
		if err != nil {
			return nil, err
		}
		emailVO = e
	}

	passVO := user.PasswordFromHash(m.Password)

	return user.RestoreUser(
		m.ID,
		m.Name,
		m.Phone,
		emailVO,
		passVO,
		m.CreatedAt,
	), nil
}

type scanner interface {
	Scan(dest ...any) error
}

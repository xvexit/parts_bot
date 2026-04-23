package repository

import (
	"context"

	dPayment "partsBot/internal/domain/payment"
	"partsBot/internal/infrastructure/db"
)

type PostgresPaymentRepo struct {
	db *db.DB
}

func NewPostgresPaymentRepo(db *db.DB) *PostgresPaymentRepo {
	return &PostgresPaymentRepo{
		db: db,
	}
}

func (r *PostgresPaymentRepo) Create(ctx context.Context, p *dPayment.Payment) (*dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		INSERT INTO payments (
			order_id,
			amount,
			provider,
			provider_txn_id,
			payment_url,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int64

	err := exec.QueryRow(
		ctx,
		query,
		p.OrderID(),
		p.Amount(),
		p.Provider(),
		p.ProviderTxnID(),
		p.PaymentURL(),
		p.Status(),
		p.CreatedAt(),
		p.UpdatedAt(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	p.SetID(id)
	return p, nil
}

func (r *PostgresPaymentRepo) Update(ctx context.Context, p *dPayment.Payment) (*dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		UPDATE payments
		SET
			amount = $1,
			provider = $2,
			provider_txn_id = $3,
			payment_url = $4,
			status = $5,
			updated_at = $6
		WHERE id = $7
	`

	_, err := exec.Exec(
		ctx,
		query,
		p.Amount(),
		p.Provider(),
		p.ProviderTxnID(),
		p.PaymentURL(),
		p.Status(),
		p.UpdatedAt(),
		p.ID(),
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PostgresPaymentRepo) GetByID(ctx context.Context, id int64) (*dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		SELECT
			id,
			order_id,
			amount,
			provider,
			provider_txn_id,
			payment_url,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	dto := PaymentModel{}

	err := exec.QueryRow(ctx, query, id).Scan(
		&dto.ID,
		&dto.OrderID,
		&dto.Amount,
		&dto.Provider,
		&dto.ProviderTxnID,
		&dto.PaymentURL,
		&dto.Status,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return mapPaymentDTOToDomain(dto), nil
}

func (r *PostgresPaymentRepo) GetByOrderID(ctx context.Context, orderID int64) (*dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		SELECT
			id,
			order_id,
			amount,
			provider,
			provider_txn_id,
			payment_url,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`

	dto := PaymentModel{}

	err := exec.QueryRow(ctx, query, orderID).Scan(
		&dto.ID,
		&dto.OrderID,
		&dto.Amount,
		&dto.Provider,
		&dto.ProviderTxnID,
		&dto.PaymentURL,
		&dto.Status,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return mapPaymentDTOToDomain(dto), nil
}

func (r *PostgresPaymentRepo) GetByProviderTxnID(ctx context.Context, providerTxnID string) (*dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		SELECT
			id,
			order_id,
			amount,
			provider,
			provider_txn_id,
			payment_url,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE provider_txn_id = $1
	`

	dto := PaymentModel{}

	err := exec.QueryRow(ctx, query, providerTxnID).Scan(
		&dto.ID,
		&dto.OrderID,
		&dto.Amount,
		&dto.Provider,
		&dto.ProviderTxnID,
		&dto.PaymentURL,
		&dto.Status,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return mapPaymentDTOToDomain(dto), nil
}

func (r *PostgresPaymentRepo) ListByOrderID(ctx context.Context, orderID int64) ([]dPayment.Payment, error) {
	exec := r.db.Executor(ctx)

	query := `
		SELECT
			id,
			order_id,
			amount,
			provider,
			provider_txn_id,
			payment_url,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC, id DESC
	`

	rows, err := exec.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]dPayment.Payment, 0)

	for rows.Next() {
		dto := PaymentModel{}

		err := rows.Scan(
			&dto.ID,
			&dto.OrderID,
			&dto.Amount,
			&dto.Provider,
			&dto.ProviderTxnID,
			&dto.PaymentURL,
			&dto.Status,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, *mapPaymentDTOToDomain(dto))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func mapPaymentDTOToDomain(dto PaymentModel) *dPayment.Payment {
	return dPayment.RestorePayment(
		dto.ID,
		dto.OrderID,
		dto.Amount,
		dto.Provider,
		dto.ProviderTxnID,
		dto.PaymentURL,
		dto.Status,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

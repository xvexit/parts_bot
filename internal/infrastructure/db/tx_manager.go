package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (m *TxManager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	err = fn(ctx)

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
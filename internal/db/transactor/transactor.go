package transactor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor struct {
	pgxPool *pgxpool.Pool
}

type txKey struct{}

func NewTransactor(pgxPool *pgxpool.Pool) *Transactor {
	return &Transactor{
		pgxPool: pgxPool,
	}
}

func (t *Transactor) CreateTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("db.Transactor.CreateTransaction: cannot begin transaction: %w", err)
	}

	defer func() {
		switch p := recover(); {
		case p != nil:
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: tx execute with panic:: transaction panic: %v, rollback err: %v", p, rbErr)
			}
		case err != nil:
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: transaction err: %v, rollback err: %v", err, rbErr)
			}
		default:
			if err = tx.Commit(ctx); err != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: cannot commit transaction: %v", err)
			}
		}
	}()

	return fn(injectTx(ctx, tx))
}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

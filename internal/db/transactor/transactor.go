package transactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCancelTx = errors.New("transaction canceled")

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
	var err error
	tx := getExecutor(ctx)
	if tx == nil {
		tx, err = t.pgxPool.BeginTx(ctx, pgx.TxOptions{})
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

		err = fn(injectTx(ctx, tx))
		if err != nil {
			return err
		}
		return nil
	}

	err = fn(ctx)
	if err != nil {
		return err
	}
	return nil
}

func getExecutor(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}

func (t *Transactor) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	tx := getExecutor(ctx)
	if tx == nil {
		return t.pgxPool.QueryRow(ctx, sql, args...)
	}

	return tx.QueryRow(ctx, sql, args...)
}

func (t *Transactor) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	tx := getExecutor(ctx)
	if tx == nil {
		return t.pgxPool.Query(ctx, sql, args...)
	}
	return tx.Query(ctx, sql, args...)
}

func (t *Transactor) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	tx := getExecutor(ctx)
	if tx == nil {
		return t.pgxPool.Exec(ctx, sql, args...)
	}
	return tx.Exec(ctx, sql, args...)
}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

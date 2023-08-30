package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Transactor struct {
	db *sql.DB
}

type txKey struct{}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

func (t *Transactor) CreateTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("db.Transactor.CreateTransaction: cannot begin transaction: %w", err)
	}

	defer func() {
		switch p := recover(); {
		case p != nil:
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: tx execute with panic:: transaction panic: %v, rollback err: %v", p, rbErr)
			}
		case err != nil:
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: transaction err: %v, rollback err: %v", err, rbErr)
			}
		default:
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("db.Transactor.WithTransaction: cannot commit transaction: %v", err)
			}
		}
	}()

	return fn(injectTx(ctx, tx))
}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

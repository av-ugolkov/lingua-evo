package database

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"lingua-evo/storage"
)

type Database struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Database {
	return &Database{
		db: pool,
	}
}

func (d Database) AddWord(w *storage.Word) error {
	return nil
}

func (d Database) EditWord(w *storage.Word) error {
	return nil
}

func (d Database) RemoveWord(w *storage.Word) error {
	return nil
}

func (d Database) PickRandomWord(w *storage.Word) (*storage.Word, error) {
	return nil, nil
}

func (d Database) SharedWord(w *storage.Word) (*storage.Word, error) {
	return nil, nil
}

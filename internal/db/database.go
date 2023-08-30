package repository

import (
	"database/sql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

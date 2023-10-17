package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

var ErrNoSavePages = errors.New("no saved page")

func NewDB(connString string) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse DB connection string: %w", err)
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db %s: %w", connString, err)
	}
	return db, nil
}

/*func (p *Word) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Value); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (Value): %w", err)
	}
	if _, err := io.WriteString(h, string(p.UserID)); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (UserID): %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}*/

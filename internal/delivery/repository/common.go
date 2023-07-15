package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

type Word struct {
	Text     string
	Language string
}

type Language struct {
	Code string
	Lang string
}

type Example struct {
	Id        uuid.UUID
	Original  string
	Translate string
}

type Dictionary struct {
	UserId        uuid.UUID
	OriginalWord  uuid.UUID
	OriginalLang  string
	TranslateLang string
	TranslateWord []uuid.UUID
	Pronunciation string
	Example       []uuid.UUID
}

type User struct {
	Username     string
	Email        string
	PasswordHash string
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

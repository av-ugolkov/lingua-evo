package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"lingua-evo/storage"
	"time"
)

type Database struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Database {
	return &Database{
		db: pool,
	}
}

func (d *Database) AddUser(ctx context.Context, userId int, userName string) error {
	query := `insert into users (user_id, user_mame) values ($1, $2)`
	_, err := d.db.Query(ctx, query, userId, userName)
	if err != nil {
		return fmt.Errorf("database.AddUser.QueryRow: %w", err)
	}
	return nil
}

func (d *Database) AddWord(ctx context.Context, w *storage.Word) error {
	userId, err := d.userId(ctx, w.UserName)
	if err != nil {
		return fmt.Errorf("database.AddWord.userId: %w", err)
	}

	queryInsWord := `insert into words (value, translate) values($1, $2) returning id`
	wordId := 0
	err = d.db.QueryRow(ctx, queryInsWord, w.Value, w.Translate[0]).Scan(&wordId)
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	queryInsDictionary := `insert into dictionary (user_id, word_id, created) values($1, $2, $3)`

	_, err = d.db.Query(ctx, queryInsDictionary, userId, wordId, time.Now())
	if err != nil {
		return fmt.Errorf("database.AddWord.Query: %w", err)
	}
	return nil
}

func (d *Database) EditWord(ctx context.Context, w *storage.Word) error {
	return nil
}

func (d *Database) FindWord(ctx context.Context, w string) (*storage.Word, error) {
	return nil, nil
}

func (d *Database) RemoveWord(ctx context.Context, w *storage.Word) error {
	return nil
}

func (d *Database) PickRandomWord(ctx context.Context, w *storage.Word) (*storage.Word, error) {
	return nil, nil
}

func (d *Database) SharedWord(ctx context.Context, w *storage.Word) (*storage.Word, error) {
	return nil, nil
}

func (d *Database) userId(ctx context.Context, userName string) (int, error) {
	var id int
	queryUser := `select id from users where user_name = ($1)`
	err := d.db.QueryRow(ctx, queryUser, userName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("database.userId.QueryRow: not found user [%s]: %w", userName, err)
	}
	return id, nil
}

package repository

import (
	"context"
	"fmt"
)

func (d *Database) AddWord(ctx context.Context, w *Word) error {
	queryInsWord := `insert into word (text, lang, pronunciation) values($1, $2, $3) returning id`
	wordId := 0
	err := d.db.QueryRow(ctx, queryInsWord, w.Value, "lang", "pronuc").Scan(&wordId) //TODO fix parameters
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return nil
}

func (d *Database) EditWord(ctx context.Context, w *Word) error {
	return nil
}

func (d *Database) FindWord(ctx context.Context, w string) (*Word, error) {
	return nil, nil
}

func (d *Database) RemoveWord(ctx context.Context, w *Word) error {
	return nil
}

func (d *Database) PickRandomWord(ctx context.Context, w *Word) (*Word, error) {
	return nil, nil
}

func (d *Database) SharedWord(ctx context.Context, w *Word) (*Word, error) {
	return nil, nil
}

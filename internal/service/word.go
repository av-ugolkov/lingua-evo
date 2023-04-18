package service

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/delivery/repository"
)

func (l *Lingua) SendWord(ctx context.Context, origWord *repository.Word) (uuid.UUID, error) {
	wordId, err := l.db.AddWord(ctx, origWord)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = l.db.AddWordInDictionary(ctx, "", wordId, wordId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return wordId, nil
}

func (l *Lingua) EditWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (l *Lingua) FindWord(ctx context.Context, w string) (*repository.Word, error) {
	return nil, nil
}

func (l *Lingua) RemoveWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (l *Lingua) PickRandomWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}

func (l *Lingua) SharedWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}

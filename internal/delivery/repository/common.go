package repository

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavePages = errors.New("no saved page")

type Storage interface {
	AddUser(ctx context.Context, userId int, userName string) error
	AddWord(ctx context.Context, w *Word) error
	EditWord(ctx context.Context, w *Word) error
	FindWord(ctx context.Context, w string) (*Word, error)
	RemoveWord(ctx context.Context, w *Word) error
	PickRandomWord(ctx context.Context, w *Word) (*Word, error)
	SharedWord(ctx context.Context, w *Word) (*Word, error)
}

type Word struct {
	UserID    int
	Value     string
	Translate []string
	Language  Language
	Example   []Example
}

type Language struct {
	Origin    string
	Translate string
}

type Example struct {
	Value     string
	Translate string
}

func (p *Word) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Value); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (Value): %w", err)
	}
	if _, err := io.WriteString(h, string(p.UserID)); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (UserID): %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

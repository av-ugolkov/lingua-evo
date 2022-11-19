package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavePages = errors.New("no saved page")

type Storage interface {
	AddWord(w *Word) error
	EditWord(w *Word) error
	RemoveWord(w *Word) error
	PickRandomWord(w *Word) (*Word, error)
	SharedWord(w *Word) (*Word, error)
}

type Word struct {
	UserName  string
	Value     string
	Translate []string
	Language  language
	Example   []example
}

type language struct {
	Origin    string
	Translate string
}

type example struct {
	Value     string
	Translate string
}

func (p *Word) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Value); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (Value): %w", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("storage.Hash.WriteString (UserName): %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

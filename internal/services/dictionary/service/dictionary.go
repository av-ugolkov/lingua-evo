package service

import (
	"context"

	"github.com/gofrs/uuid"
)

type (
	repoDict interface {
		AddDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, error)
		DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error
	}
)

type DictionarySvc struct {
	repo repoDict
}

func NewDictionarySvc(repo repoDict) *DictionarySvc {
	return &DictionarySvc{
		repo: repo,
	}
}

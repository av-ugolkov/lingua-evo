package service

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	runtimeAccess "github.com/av-ugolkov/lingua-evo/runtime/access"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

type (
	repoVocabAccess interface {
		AddAccessForUser(ctx context.Context, vid, uid uuid.UUID, isEditor bool) error
		RemoveAccessForUser(ctx context.Context, vid, uid uuid.UUID) error
		GetEditable(ctx context.Context, vid, uid uuid.UUID) (bool, error)
	}
)

func (s *Service) VocabularyEditable(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	editable, err := s.repoVocab.GetEditable(ctx, vid, uid)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("vocabulary_access.Service.VocabularyEditable: %w", err)
	default:
		return editable, nil
	}
}

func (s *Service) AddAccessForUser(ctx context.Context, access entity.Access) error {
	err := s.repoVocab.AddAccessForUser(ctx, access.VocabID, access.UserID, access.Status == runtimeAccess.Edit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.AddAccessForUser: %w", err)
	}
	return nil
}

func (s *Service) RemoveAccessForUser(ctx context.Context, access entity.Access) error {
	err := s.repoVocab.RemoveAccessForUser(ctx, access.VocabID, access.UserID)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.RemoveAccessForUser: %w", err)
	}
	return nil
}

func (s *Service) UpdateAccessForUser(ctx context.Context, access entity.Access) error {
	err := s.RemoveAccessForUser(ctx, access)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.UpdateAccessForUser - remove: %w", err)
	}

	err = s.AddAccessForUser(ctx, access)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.UpdateAccessForUser - add: %w", err)
	}

	return nil
}

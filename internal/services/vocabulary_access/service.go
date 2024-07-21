package vocabulary_access

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	repoVocabAccess interface {
		AddAccessForUser(ctx context.Context, vocabID, userID uuid.UUID, isEditor bool) error
		RemoveAccessForUser(ctx context.Context, vocabID, userID uuid.UUID) error
		GetAccess(ctx context.Context, vocabID, userID uuid.UUID) (bool, error)
	}
)

type Service struct {
	repoVocabAccess repoVocabAccess
}

func NewService(repoVocabAccess repoVocabAccess) *Service {
	return &Service{
		repoVocabAccess: repoVocabAccess,
	}
}

func (s *Service) VocabularyEditable(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	editable, err := s.repoVocabAccess.GetAccess(ctx, vid, uid)
	if err != nil {
		return false, fmt.Errorf("vocabulary_access.Service.GetVocabularyAccess: %w", err)
	}
	return editable, nil
}

func (s *Service) AddAccessForUser(ctx context.Context, access Access) error {
	err := s.repoVocabAccess.AddAccessForUser(ctx, access.VocabID, access.UserID, access.AccessEdit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.AddAccessForUser: %w", err)
	}
	return nil
}

func (s *Service) RemoveAccessForUser(ctx context.Context, access Access) error {
	err := s.repoVocabAccess.RemoveAccessForUser(ctx, access.VocabID, access.UserID)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.RemoveAccessForUser: %w", err)
	}
	return nil
}

func (s *Service) UpdateAccessForUser(ctx context.Context, access Access) error {
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

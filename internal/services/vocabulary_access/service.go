package vocabulary_access

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	repoVocab interface {
		ChangeAccess(ctx context.Context, vocabID uuid.UUID, access int, accessEdit bool) error
		AddAccessForUser(ctx context.Context, vocabID, userID uuid.UUID, isEditor bool) error
		RemoveAccessForUser(ctx context.Context, vocabID, userID uuid.UUID) error
	}
)

type Service struct {
	repoVocab repoVocab
}

func NewService(repoVocab repoVocab) *Service {
	return &Service{
		repoVocab: repoVocab,
	}
}

func (s *Service) ChangeAccess(ctx context.Context, access Access) error {
	err := s.repoVocab.ChangeAccess(ctx, access.VocabID, access.ID, access.AccessEdit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.ChangeAccess: %w", err)
	}
	return nil
}

func (s *Service) GetVocabularyAccess(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	return false, nil
}

func (s *Service) AddAccessForUser(ctx context.Context, access Access) error {
	err := s.repoVocab.AddAccessForUser(ctx, access.VocabID, access.UserID, access.AccessEdit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.AddAccessForUser: %w", err)
	}
	return nil
}

func (s *Service) RemoveAccessForUser(ctx context.Context, access Access) error {
	err := s.repoVocab.RemoveAccessForUser(ctx, access.VocabID, access.UserID)
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

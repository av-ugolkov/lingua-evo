package vocabulary_access

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"

	"github.com/google/uuid"
)

type (
	repoVocab interface {
		ChangeAccess(ctx context.Context, vocabID uuid.UUID, access int, accessEdit bool) error
		AddAccessForUser(ctx context.Context, vocabID, userID uuid.UUID, isEditor bool) error
		RemoveAccessForUser(ctx context.Context, vocabID, userID uuid.UUID) error
	}

	vocabSvc interface {
		IsVocabularyOwner(ctx context.Context, vocabID uuid.UUID, userID uuid.UUID) (bool, error)
	}
)

type Service struct {
	tr        *transactor.Transactor
	repoVocab repoVocab
	vocabSvc  vocabSvc
}

func NewService(tr *transactor.Transactor, repoVocab repoVocab, vocabSvc vocabSvc) *Service {
	return &Service{
		tr:        tr,
		repoVocab: repoVocab,
		vocabSvc:  vocabSvc,
	}
}

func (s *Service) ChangeAccess(ctx context.Context, access Access) error {
	err := s.repoVocab.ChangeAccess(ctx, access.VocabID, access.ID, access.AccessEdit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.Service.ChangeAccess: %w", err)
	}
	return nil
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

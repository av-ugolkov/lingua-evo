package vocab_words

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	vocabSvc interface {
		CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error)
	}

	wordSvc interface {
		CopyWords(ctx context.Context, vid, copyVid uuid.UUID) error
	}
)

type Service struct {
	vocabSvc vocabSvc
	wordSvc  wordSvc
}

func NewService(vocabSvc vocabSvc, wordSvc wordSvc) *Service {
	return &Service{
		vocabSvc: vocabSvc,
		wordSvc:  wordSvc,
	}
}

func (s *Service) CopyVocab(ctx context.Context, uid, vid uuid.UUID) error {
	copyVid, err := s.vocabSvc.CopyVocab(ctx, uid, vid)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.Copy - copy vocabulary: %w", err)
	}

	err = s.wordSvc.CopyWords(ctx, vid, copyVid)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.Copy - copy words: %w", err)
	}
	return nil
}

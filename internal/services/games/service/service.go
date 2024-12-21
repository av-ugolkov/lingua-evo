package service

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/games"

	"github.com/google/uuid"
)

type (
	repo interface {
		GerWords(ctx context.Context, uid, vid uuid.UUID, count int) ([]entity.ReviseGameWord, error)
	}
)

type Service struct {
	repo repo
}

func New(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GameRevise(ctx context.Context, uid uuid.UUID, data entity.Game) error {
	if data.TypeGame != entity.TypeGameRevise {
		return msgerr.New(fmt.Errorf("games.Service.GameRevise: game is not revise"), entity.ErrMsgWrongDataGame)
	}

	words, err := s.repo.GerWords(ctx, uid, data.VocabID, data.CountWord)
	if err != nil {
		return msgerr.New(err, msgerr.ErrMsgInternal)
	}

	fmt.Println(words)

	return nil
}

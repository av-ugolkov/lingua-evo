package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

//go:generate mockery --inpackage --outpkg notifications --testonly --name "repoNotification"

type (
	repoNotification interface {
		GetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error)
		SetVocabNotification(ctx context.Context, uid, vid uuid.UUID) error
		DeleteVocabNotification(ctx context.Context, uid, vid uuid.UUID) error
		GetVocabNotifications(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
	}
)

type Service struct {
	repo repoNotification
}

func NewService(repo repoNotification) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	ok, err := s.repo.GetVocabNotification(ctx, uid, vid)
	if err != nil {
		return false, fmt.Errorf("notifications.Service.GetVocabNotification: %w", err)
	}

	return ok, nil
}

func (s *Service) SetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	ok, err := s.GetVocabNotification(ctx, uid, vid)
	if err != nil {
		return false, fmt.Errorf("notifications.Service.SetVocabNotifications: %w", err)
	}
	if ok {
		err = s.repo.DeleteVocabNotification(ctx, uid, vid)
		if err != nil {
			return false, fmt.Errorf("notifications.Service.SetVocabNotifications: %w", err)
		}
		return false, nil
	}
	err = s.repo.SetVocabNotification(ctx, uid, vid)
	if err != nil {
		return false, fmt.Errorf("notifications.Service.SetVocabNotifications: %w", err)
	}
	return true, nil
}

func (s *Service) GetVocabNotifications(ctx context.Context, uid uuid.UUID) (_ []uuid.UUID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("notifications.Service.GetVocabNotifications: %w", err)
		}
	}()

	vocabNotifications, err := s.repo.GetVocabNotifications(ctx, uid)
	if err != nil {
		return nil, err
	}

	return vocabNotifications, nil
}

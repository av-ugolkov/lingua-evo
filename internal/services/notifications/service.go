package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	notificationRepo interface {
		GetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error)
		SetVocabNotification(ctx context.Context, uid, vid uuid.UUID) error
		DeleteVocabNotification(ctx context.Context, uid, vid uuid.UUID) error
	}
)

type Service struct {
	repo notificationRepo
}

func NewService(repo notificationRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	ok, err := s.repo.GetVocabNotification(ctx, uid, vid)
	if err != nil {
		return false, fmt.Errorf("notifications.Service.GetVocabNotifications: %w", err)
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

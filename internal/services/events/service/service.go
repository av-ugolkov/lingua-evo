package service

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"

	"github.com/google/uuid"
)

type (
	repoEvent interface {
		GetCountEvents(ctx context.Context, subscriberIDs []uuid.UUID) (int, error)
		GetEvents(ctx context.Context, subscriberIDs []uuid.UUID) ([]entity.Event, error)
	}

	notificationsSvc interface {
		GetVocabNotifications(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
	}
)

type Service struct {
	tr               *transactor.Transactor
	repoEvent        repoEvent
	notificationsSvc notificationsSvc
}

func NewService(
	tr *transactor.Transactor,
	repoEvent repoEvent,
	notificationsSvc notificationsSvc,
) *Service {
	return &Service{
		tr:               tr,
		repoEvent:        repoEvent,
		notificationsSvc: notificationsSvc,
	}
}

func (s *Service) GetCountEvents(ctx context.Context, uid uuid.UUID) (_ int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.Service.GetCountEvents: %w", err)
		}
	}()

	vocabIDs, err := s.notificationsSvc.GetVocabNotifications(ctx, uid)
	if err != nil {
		return 0, err
	}

	if len(vocabIDs) == 0 {
		return 0, nil
	}

	count, err := s.repoEvent.GetCountEvents(ctx, vocabIDs)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (s *Service) GetEvents(ctx context.Context, uid uuid.UUID) (_ []entity.Event, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.Service.GetEvents: %w", err)
		}
	}()

	vocabIDs, err := s.notificationsSvc.GetVocabNotifications(ctx, uid)
	if err != nil {
		return nil, err
	}

	events, err := s.repoEvent.GetEvents(ctx, vocabIDs)
	if err != nil {
		return nil, err
	}

	return events, nil
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"

	"github.com/google/uuid"
)

type (
	repoEvents interface {
		GetCountVocabEvents(ctx context.Context, subscriberIDs []uuid.UUID) (int, error)
		GetVocabEvents(ctx context.Context, subscriberIDs []uuid.UUID) ([]entity.Event, error)
		AddEvent(ctx context.Context, uid uuid.UUID, payload entity.Payload) error
	}

	notificationsSvc interface {
		GetVocabNotifications(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
	}
)

//go:generate mockery --inpackage --outpkg events --testonly --name "notificationsSvc"

type Service struct {
	tr               *transactor.Transactor
	repoEvents       repoEvents
	notificationsSvc notificationsSvc
}

func NewService(
	tr *transactor.Transactor,
	repoEvents repoEvents,
	notificationsSvc notificationsSvc,
) *Service {
	return &Service{
		tr:               tr,
		repoEvents:       repoEvents,
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

	count, err := s.repoEvents.GetCountVocabEvents(ctx, vocabIDs)
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

	events, err := s.repoEvents.GetVocabEvents(ctx, vocabIDs)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) AddEvent(ctx context.Context, uid uuid.UUID, payload entity.Payload) error {
	err := s.repoEvents.AddEvent(ctx, uid, payload)
	if err != nil {
		return fmt.Errorf("events.Service.AddEvent: %w", err)
	}
	return nil
}

func (s *Service) AsyncAddEvent(uid uuid.UUID, payload entity.Payload) error {
	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFn()

	return s.AddEvent(timeoutCtx, uid, payload)
}

func (s *Service) ReadEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) error {
	return nil
}

func (s *Service) UpdateEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) error {
	return nil
}

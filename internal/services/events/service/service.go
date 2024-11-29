package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type (
	repoEvents interface {
		GetCountVocabEvents(ctx context.Context, subscriberIDs []uuid.UUID) (int, error)
		GetVocabEvents(ctx context.Context, subscriberIDs []uuid.UUID) ([]entity.Event, error)
		AddEvent(ctx context.Context, uid uuid.UUID, typeEvent entity.PayloadType, payload []byte) (uuid.UUID, error)
		ReadEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) error
		GetWatchedEvents(ctx context.Context, uid uuid.UUID) ([]entity.EventWatched, error)
		DeleteWatchedEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) error
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

func (s *Service) AddEvent(ctx context.Context, event entity.Event) (uuid.UUID, error) {
	data, err := jsoniter.Marshal(event.Payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("events.Service.AddEvent: %w", err)
	}

	eid, err := s.repoEvents.AddEvent(ctx, event.User.ID, event.Type, data)
	if err != nil {
		return uuid.Nil, fmt.Errorf("events.Service.AddEvent: %w", err)
	}
	return eid, nil
}

func (s *Service) AsyncAddEvent(event entity.Event) {
	go func() {
		timeoutCtx, cancelFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFn()

		_, err := s.AddEvent(timeoutCtx, event)
		if err != nil {
			slog.Error(fmt.Sprintf("events.Service.AsyncAddEvent: %v", err))
		}
	}()
}

func (s *Service) ReadEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) error {
	err := s.repoEvents.ReadEvent(ctx, uid, eventID)
	if err != nil {
		return fmt.Errorf("events.Service.ReadEvent: %w", err)
	}

	return nil
}

package subscribers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	repoSubscribers interface {
		Get(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
		GetRespondents(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
		Subscribe(ctx context.Context, uid, subID uuid.UUID) error
		Unsubscribe(ctx context.Context, uid, subID uuid.UUID) error
		Check(ctx context.Context, uid, subID uuid.UUID) (bool, error)
	}
)

type Service struct {
	repoSubscribers repoSubscribers
}

func NewService(repoSubscribers repoSubscribers) *Service {
	return &Service{
		repoSubscribers: repoSubscribers,
	}
}

func (s *Service) Get(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	subscribers, err := s.repoSubscribers.Get(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("subscribers.Service.GetSubscribers: %w", err)
	}

	return subscribers, nil
}

func (s *Service) GetRespondents(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	subscribers, err := s.repoSubscribers.GetRespondents(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("subscribers.Service.GetSubscribers: %w", err)
	}

	return subscribers, nil
}

func (s *Service) Subscribe(ctx context.Context, uid, subID uuid.UUID) error {
	err := s.repoSubscribers.Subscribe(ctx, uid, subID)
	if err != nil {
		return fmt.Errorf("subscribers.Service.Subscribe: %w", err)
	}
	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, uid, subID uuid.UUID) error {
	err := s.repoSubscribers.Unsubscribe(ctx, uid, subID)
	if err != nil {
		return fmt.Errorf("subscribers.Service.Unsubscribe: %w", err)
	}
	return nil
}

func (s *Service) Check(ctx context.Context, uid, subID uuid.UUID) (bool, error) {
	isSubscriber, err := s.repoSubscribers.Check(ctx, uid, subID)
	if err != nil {
		return false, fmt.Errorf("subscribers.Service.Check: %w", err)
	}
	return isSubscriber, nil
}

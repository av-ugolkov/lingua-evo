package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"

	"github.com/google/uuid"
)

type EventRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *EventRepo {
	return &EventRepo{
		tr: tr,
	}
}

func (r *EventRepo) GetCountEvents(ctx context.Context, vocabIDs []uuid.UUID) (int, error) {
	const query = `SELECT COUNT(id) FROM event_vocab WHERE vocab_id=ANY($1);`

	var count int
	err := r.tr.QueryRow(ctx, query, vocabIDs).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("events.delivery.repository.UserRepo.GetCountUserEvents: %w", err)
	}

	return count, nil
}

func (r *EventRepo) GetEvents(ctx context.Context, vocabIDs []uuid.UUID) ([]entity.Event, error) {
	return nil, nil
}

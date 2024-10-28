package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *EventRepo) GetCountVocabEvents(ctx context.Context, vocabIDs []uuid.UUID) (int, error) {
	const query = `
		SELECT COUNT(id) FROM events e
		LEFT JOIN vocabulary_notifications vn ON vn.vocab_id::text = e.payload->'data'->>'VocabID'
		WHERE payload->'data'->>'VocabID'=ANY($1)
		AND e.created_at >= vn.created_at;`

	var count int
	err := r.tr.QueryRow(ctx, query, vocabIDs).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("events.delivery.repository.UserRepo.GetCountUserEvents: %w", err)
	}

	return count, nil
}

func (r *EventRepo) GetVocabEvents(ctx context.Context, vocabIDs []uuid.UUID) ([]entity.Event, error) {
	const query = `
		SELECT e.id, e.user_id, e.payload, e.created_at FROM events e
		LEFT JOIN vocabulary_notifications vn ON vn.vocab_id::text = e.payload->'data'->>'VocabID'
		WHERE payload->'data'->>'VocabID'=ANY($1)
		AND e.created_at >= vn.created_at;`

	rows, err := r.tr.Query(ctx, query, vocabIDs)
	if err != nil {
		return nil, fmt.Errorf("events.delivery.repository.UserRepo.GetEventsVocab: %w", err)
	}

	events := make([]entity.Event, 0, 10)
	if rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.ID, &event.UserID, &event.Payload, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("events.delivery.repository.UserRepo.GetEventsVocab: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepo) AddEvent(ctx context.Context, uid uuid.UUID, payload entity.Payload) error {
	const query = `
		INSERT INTO events (id, user_id, payload, created_at) 
		VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, payload) DO UPDATE SET created_at = $4;`

	_, err := r.tr.Exec(ctx, query, uuid.New(), uid, payload, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("events.delivery.repository.UserRepo.AddEvent: %w", err)
	}

	return nil
}

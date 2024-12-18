package repository

import (
	"context"
	"fmt"
	"log/slog"
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
		SELECT COUNT(e.id)-COUNT(ew.event_id) FROM events e
		LEFT JOIN vocabulary_notifications vn ON vn.vocab_id::text = e.payload->>'vocab_id'
		LEFT JOIN events_watched ew ON ew.event_id = e.id
		WHERE payload->>'vocab_id'=ANY($1)
		AND e.created_at >= vn.created_at;`

	var count int
	err := r.tr.QueryRow(ctx, query, vocabIDs).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("events.repository.UserRepo.GetCountUserEvents: %w", err)
	}

	return count, nil
}

func (r *EventRepo) GetVocabEvents(ctx context.Context, vocabIDs []uuid.UUID) ([]entity.Event, error) {
	const query = `
		SELECT e.id, e.user_id, u.nickname, et.name, e.payload, COALESCE(e.id=ew.event_id, false), e.created_at FROM events e
		LEFT JOIN vocabulary_notifications vn ON vn.vocab_id::text = e.payload->>'vocab_id'
		LEFT JOIN events_type et ON et.id = e.type
		LEFT JOIN users u ON u.id = e.user_id
		LEFT JOIN events_watched ew ON ew.event_id = e.id
		WHERE payload->>'vocab_id'=ANY($1)
		AND e.created_at >= vn.created_at
		ORDER BY e.created_at DESC;`

	rows, err := r.tr.Query(ctx, query, vocabIDs)
	if err != nil {
		return nil, fmt.Errorf("events.repository.UserRepo.GetEventsVocab: %w", err)
	}

	events := make([]entity.Event, 0, 10)
	for rows.Next() {
		var event entity.Event
		var jsonData []byte
		if err := rows.Scan(&event.ID, &event.User.ID, &event.User.Nickname, &event.Type, &jsonData, &event.Watched, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("events.repository.UserRepo.GetEventsVocab: %w", err)
		}

		event.Payload, err = entity.Unmarshal(event.Type, jsonData)
		if err != nil {
			slog.Error(fmt.Sprintf("events.repository.UserRepo.GetEventsVocab: %v", err))
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepo) AddEvent(ctx context.Context, uid uuid.UUID, typeEvent entity.PayloadType, payload []byte) (uuid.UUID, error) {
	const query = `
		INSERT INTO events (id, user_id, type, payload, created_at) 
		VALUES ($1, $2, (SELECT id FROM events_type WHERE "name"=$3), $4, $5) 
		ON CONFLICT (user_id, type, payload) 
		DO UPDATE SET created_at = $5 
		RETURNING id;`

	var eid uuid.UUID
	err := r.tr.QueryRow(ctx, query, uuid.New(), uid, typeEvent, payload, time.Now().UTC()).Scan(&eid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("events.repository.UserRepo.AddEvent: %w", err)
	}

	return eid, nil
}

func (r *EventRepo) ReadEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.repository.UserRepo.ReadEvent: %w", err)
		}
	}()

	const query = `
		INSERT INTO events_watched (event_id, user_id, watched_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) DO NOTHING;`

	result, err := r.tr.Exec(ctx, query, eventID, uid, time.Now().UTC())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("change 0 or more than 1 rows")
	}

	return nil
}

func (r *EventRepo) GetWatchedEvents(ctx context.Context, uid uuid.UUID) ([]entity.EventWatched, error) {
	const query = `
		SELECT w.event_id, w.user_id, w.watched_at 
		FROM events_watched w
		WHERE w.user_id = $1;`

	rows, err := r.tr.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("events.repository.UserRepo.GetWatchedEvents: %w", err)
	}

	eventsWatched := make([]entity.EventWatched, 0, 10)
	for rows.Next() {
		var event entity.EventWatched
		if err := rows.Scan(&event.EventID, &event.UserID, &event.WatchedAt); err != nil {
			return nil, fmt.Errorf("events.repository.UserRepo.GetWatchedEvents: %w", err)
		}
		eventsWatched = append(eventsWatched, event)
	}

	return eventsWatched, nil
}

func (r *EventRepo) DeleteWatchedEvent(ctx context.Context, uid uuid.UUID, eventID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.repository.UserRepo.DeleteWatchedEvent: %w", err)
		}
	}()

	const query = `
		DELETE FROM events_watched WHERE event_id = $1 AND user_id = $2;`

	result, err := r.tr.Exec(ctx, query, eventID, uid)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("change 0 or more than 1 rows")
	}

	return nil
}

package telegram

import (
	"context"
	"errors"
	"fmt"

	clientsTelegram "lingua-evo/internal/clients/telegram"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/internal/events"
	eventsCommands "lingua-evo/internal/events/commands"
)

type Processor struct {
	tg       *clientsTelegram.Client
	database repository.Database
	lastCmd  eventsCommands.Command
	offset   int
}

type Meta struct {
	ChatID   int
	UserID   int
	UserName string
	IsBot    bool
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *clientsTelegram.Client, database repository.Database) *Processor {
	return &Processor{
		tg:       client,
		database: database,
		lastCmd:  eventsCommands.UnknownCmd,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("telegram.Fetch.Updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("telegram.Process: %w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("telegram.processMessage.meta: %w", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserID, meta.UserName); err != nil {
		return fmt.Errorf("telegram.processMessage.doCmd: %w", err)
	}

	return nil
}

func (p *Processor) sendStart(chatID int, userId int, userName string) error {
	_, err := p.database.AddUser(context.Background(), &repository.User{Username: userName})
	if err != nil {
		return fmt.Errorf("telegram.sendStart.AddUser: %w", err)
	}

	return p.tg.SendMessage(chatID, fmt.Sprintf(eventsCommands.MsgHello, userName))
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, eventsCommands.MsgHelp)
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("telegram.meta: %w", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd clientsTelegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {
	case events.Message:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserID:   upd.Message.From.ID,
			IsBot:    upd.Message.From.IsBot,
			UserName: upd.Message.From.Username,
		}
	case events.CallbackQuery:
		res.Meta = Meta{
			ChatID:   upd.CallbackQuery.Chat.ID,
			UserID:   upd.CallbackQuery.From.ID,
			IsBot:    upd.CallbackQuery.From.IsBot,
			UserName: upd.CallbackQuery.From.Username,
		}
	}

	return res
}

func fetchType(upd clientsTelegram.Update) events.Type {
	switch {
	case upd.Message != nil:
		return events.Message
	case upd.CallbackQuery != nil:
		return events.CallbackQuery
	default:
		return events.Unknown
	}
}

func fetchText(upd clientsTelegram.Update) string {
	switch {
	case upd.Message != nil:
		return upd.Message.Text
	case upd.CallbackQuery != nil:
		return upd.CallbackQuery.Data
	default:
		return ""
	}
}

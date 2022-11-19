package messengers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"lingua-evo/clients/telegram"
	"lingua-evo/events"
	"lingua-evo/events/commands"
	"lingua-evo/storage"
)

type Processor struct {
	tg      *telegram.Client
	storage storage.Storage
	lastCmd commands.Command
	offset  int
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		lastCmd: commands.UnknownCmd,
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
	/*meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("telegram.processMessage.meta: %w", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return fmt.Errorf("telegram.processMessage.doCmd: %w", err)
	}*/

	return nil
}

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from %s", text, username)
	/*cmd := p.chooseCmd(text)
	if cmd == UnknownCmd {
		if p.lastCmd != UnknownCmd {
			cmd = p.lastCmd
		}
	} else {
		p.lastCmd = cmd
	}
	switch cmd {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case AddCmd:
		return p.sendAdd(chatID, text, username)
	case Cancel:
		return p.sendCancel(chatID)
	case RndCmd:
		return p.sendRandom(chatID, username)
	default:
		return p.tg.SendMessage(chatID, telegram.msgUnknownCommand, "")
	}*/
	return nil
}

func (p *Processor) chooseCmd(cmd string) commands.Command {
	switch cmd {
	case commands.StartCmd:
		return commands.StartCmd
	case commands.HelpCmd:
		return commands.HelpCmd
	case commands.Cancel:
		return commands.Cancel
	case commands.AddCmd:
		return commands.AddCmd
	case commands.RndCmd:
		return commands.RndCmd
	default:
		return commands.UnknownCmd
	}
}

func (p *Processor) sendCancel(chatID int) error {
	p.lastCmd = commands.UnknownCmd
	return p.tg.SendMessage(chatID, commands.MsgCancel, "")
}

func (p *Processor) sendAdd(chatID int, world string, username string) error {
	return nil
}

func (p *Processor) sendRandom(chatID int, world *storage.Word) error {
	page, err := p.storage.PickRandomWord(world)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return fmt.Errorf("messengers.telegram.senRandom.PickRandom: %w", err)
	}

	if errors.Is(err, storage.ErrNoSavePages) {
		return p.tg.SendMessage(chatID, commands.MsgNoSavedPages, "")
	}

	if err := p.tg.SendMessage(chatID, page.Value, ""); err != nil {
		return fmt.Errorf("messengers.telegram.senRandom.SendMessage: %w", err)
	}

	return nil //p.storage.RemoveWord(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, commands.MsgHelp, "")
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, commands.MsgHello, "")
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("telegram.meta: %w", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

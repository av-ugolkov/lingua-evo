package telegram

import (
	"errors"
	"fmt"

	"lingua-evo/events/commands"
	"lingua-evo/storage"
)

func (p *Processor) sendRandom(chatID int, world *storage.Word) error {
	page, err := p.storage.PickRandomWord(world)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return fmt.Errorf("messengers.telegram.senRandom.PickRandom: %w", err)
	}

	if errors.Is(err, storage.ErrNoSavePages) {
		return p.tg.SendMessage(chatID, commands.MsgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.Value); err != nil {
		return fmt.Errorf("messengers.telegram.senRandom.SendMessage: %w", err)
	}

	return nil //p.storage.RemoveWord(page)
}

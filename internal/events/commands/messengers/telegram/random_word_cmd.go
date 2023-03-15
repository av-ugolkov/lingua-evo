package telegram

import (
	"context"
	"errors"
	"fmt"

	storage "lingua-evo/internal/delivery/repository"
	"lingua-evo/internal/events/commands"
)

func (p *Processor) sendRandom(chatID int, world *storage.Word) error {
	word, err := p.database.PickRandomWord(context.Background(), world)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return fmt.Errorf("messengers.telegram.senRandom.PickRandom: %w", err)
	}

	if errors.Is(err, storage.ErrNoSavePages) {
		return p.tg.SendMessage(chatID, commands.MsgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, word.Text); err != nil {
		return fmt.Errorf("messengers.telegram.senRandom.SendMessage: %w", err)
	}

	return nil //p.storage.RemoveWord(page)
}

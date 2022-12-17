package telegram

import (
	commands2 "lingua-evo/internal/events/commands"
)

func (p *Processor) sendCancel(chatID int) error {
	p.lastCmd = commands2.UnknownCmd
	return p.tg.SendMessage(chatID, commands2.MsgCancel)
}

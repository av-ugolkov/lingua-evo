package telegram

import "lingua-evo/events/commands"

func (p *Processor) sendCancel(chatID int) error {
	p.lastCmd = commands.UnknownCmd
	return p.tg.SendMessage(chatID, commands.MsgCancel)
}

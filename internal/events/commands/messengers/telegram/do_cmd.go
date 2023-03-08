package telegram

import (
	"log"
	"strings"

	storage "lingua-evo/internal/delivery/repository"
	eventsCommands "lingua-evo/internal/events/commands"
)

func (p *Processor) doCmd(text string, chatID int, userId int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from [%s]", text, userName)
	cmd := p.chooseCmd(eventsCommands.Command(text))
	if cmd == eventsCommands.UnknownCmd {
		if p.lastCmd != eventsCommands.UnknownCmd {
			cmd = p.lastCmd
		}
	} else {
		p.lastCmd = cmd
	}
	switch cmd {
	case eventsCommands.StartCmd:
		return p.sendStart(chatID, userId, userName)
	case eventsCommands.HelpCmd:
		return p.sendHelp(chatID)
	case eventsCommands.AddCmd:
		return p.addWord(chatID)
	case eventsCommands.Cancel:
		return p.sendCancel(chatID)
	case eventsCommands.RndCmd:
		return p.sendRandom(chatID, &storage.Word{})
	default:
		return p.tg.SendMessage(chatID, eventsCommands.MsgUnknownCommand)
	}
}

func (p *Processor) chooseCmd(cmd eventsCommands.Command) eventsCommands.Command {
	switch cmd {
	case eventsCommands.StartCmd:
		return eventsCommands.StartCmd
	case eventsCommands.HelpCmd:
		return eventsCommands.HelpCmd
	case eventsCommands.Cancel:
		return eventsCommands.Cancel
	case eventsCommands.AddCmd:
		return eventsCommands.AddCmd
	case eventsCommands.RndCmd:
		return eventsCommands.RndCmd
	default:
		return eventsCommands.UnknownCmd
	}
}

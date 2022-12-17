package telegram

import (
	commands2 "lingua-evo/internal/events/commands"
	"lingua-evo/pkg/storage"
	"log"
	"strings"
)

func (p *Processor) doCmd(text string, chatID int, userId int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from [%s]", text, userName)
	cmd := p.chooseCmd(commands2.Command(text))
	if cmd == commands2.UnknownCmd {
		if p.lastCmd != commands2.UnknownCmd {
			cmd = p.lastCmd
		}
	} else {
		p.lastCmd = cmd
	}
	switch cmd {
	case commands2.StartCmd:
		return p.sendStart(chatID, userId, userName)
	case commands2.HelpCmd:
		return p.sendHelp(chatID)
	case commands2.AddCmd:
		return p.addWord(chatID)
	case commands2.Cancel:
		return p.sendCancel(chatID)
	case commands2.RndCmd:
		return p.sendRandom(chatID, &storage.Word{})
	default:
		return p.tg.SendMessage(chatID, commands2.MsgUnknownCommand)
	}
}

func (p *Processor) chooseCmd(cmd commands2.Command) commands2.Command {
	switch cmd {
	case commands2.StartCmd:
		return commands2.StartCmd
	case commands2.HelpCmd:
		return commands2.HelpCmd
	case commands2.Cancel:
		return commands2.Cancel
	case commands2.AddCmd:
		return commands2.AddCmd
	case commands2.RndCmd:
		return commands2.RndCmd
	default:
		return commands2.UnknownCmd
	}
}

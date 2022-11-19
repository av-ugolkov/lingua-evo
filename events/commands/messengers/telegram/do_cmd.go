package telegram

import (
	"log"
	"strings"

	"lingua-evo/events/commands"
	"lingua-evo/storage"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from %s", text, username)
	cmd := p.chooseCmd(commands.Command(text))
	if cmd == commands.UnknownCmd {
		if p.lastCmd != commands.UnknownCmd {
			cmd = p.lastCmd
		}
	} else {
		p.lastCmd = cmd
	}
	switch cmd {
	case commands.StartCmd:
		return p.sendHello(chatID)
	case commands.HelpCmd:
		return p.sendHelp(chatID)
	case commands.AddCmd:
		return p.addWord(chatID)
	case commands.Cancel:
		return p.sendCancel(chatID)
	case commands.RndCmd:
		return p.sendRandom(chatID, &storage.Word{})
	default:
		return p.tg.SendMessage(chatID, commands.MsgUnknownCommand)
	}
}

func (p *Processor) chooseCmd(cmd commands.Command) commands.Command {
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

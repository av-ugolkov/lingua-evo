package commands

import (
	"LinguaEvo/clients/telegram"
)

var numericKeyboard = `{"inline_keyboard":[[{"text":"Some button text 1", "callback_data": "1" }],[{ "text": "Some button text 2", "callback_data": "2" }],[{ "text": "Some button text 3", "callback_data": "3" }]]}`

type AddWordCommand struct {
	tg    *telegram.Client
	index int
}

func (c *AddWordCommand) Execute(chatID int) (Command, error) {
	defer func() { c.index++ }()
	switch c.index {
	case 0:
		return AddCmd, c.tg.SendMessage(chatID, MsgAddWord, numericKeyboard)
	case 1:
		return AddCmd, c.tg.SendMessage(chatID, MsgAddPronounce, numericKeyboard)
	case 2:
		return AddCmd, c.tg.SendMessage(chatID, MsgAddTranslate, numericKeyboard)
	case 3:
		return AddCmd, c.tg.SendMessage(chatID, MsgAddExample, numericKeyboard)
	default:
		c.index = -1
		return UnknownCmd, c.tg.SendMessage(chatID, MsgAddFinish, numericKeyboard)
	}
	return UnknownCmd, nil
}

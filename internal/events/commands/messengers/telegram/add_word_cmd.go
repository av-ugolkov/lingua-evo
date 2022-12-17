package telegram

import (
	"lingua-evo/internal/clients/telegram"
	"lingua-evo/internal/events/commands"
)

const (
	word = iota
	pronounce
	translate
	example
)

var index = 0

var keyboard = telegram.InlineKeyboard{
	InlineKeyboardButton: [][]telegram.InlineKeyboardButton{
		{
			telegram.InlineKeyboardButton{Text: "RU", CallbackData: "RU"},
			telegram.InlineKeyboardButton{Text: "EN", CallbackData: "EN"},
		},
		{
			telegram.InlineKeyboardButton{Text: "FR", CallbackData: "FR"},
			telegram.InlineKeyboardButton{Text: "GE", CallbackData: "GE"},
		},
	},
}

func (p *Processor) addWord(chatID int) error {
	defer func() { index++ }()
	switch index {
	case word:
		return p.tg.SendMessage(chatID, commands.MsgAddWord)
	case pronounce:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddPronounce, keyboard.JSON())
	case translate:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddTranslate, keyboard.JSON())
	case example:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddExample, keyboard.JSON())
	default:
		index = -1
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddFinish, keyboard.JSON())
	}
}

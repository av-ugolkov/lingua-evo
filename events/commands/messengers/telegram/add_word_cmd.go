package telegram

import "lingua-evo/events/commands"

const (
	word = iota
	pronounce
	translate
	example
)

var index = 0

var numericKeyboard = `{"inline_keyboard":[[{"text":"Some button text 1", "callback_data": "1" }],[{ "text": "Some button text 2", "callback_data": "2" }],[{ "text": "Some button text 3", "callback_data": "3" }]]}`

func (p *Processor) addWord(chatID int) error {
	defer func() { index++ }()
	switch index {
	case word:
		return p.tg.SendMessage(chatID, commands.MsgAddWord)
	case pronounce:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddPronounce, numericKeyboard)
	case translate:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddTranslate, numericKeyboard)
	case example:
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddExample, numericKeyboard)
	default:
		index = -1
		return p.tg.SendMessageWithButton(chatID, commands.MsgAddFinish, numericKeyboard)
	}
}

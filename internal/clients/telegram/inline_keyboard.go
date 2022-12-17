package telegram

import (
	"encoding/json"
)

type InlineKeyboard struct {
	InlineKeyboardButton [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

func CreateInlineKeyBoard() InlineKeyboard {
	return InlineKeyboard{
		InlineKeyboardButton: [][]InlineKeyboardButton{},
	}
}

func (k *InlineKeyboard) JSON() string {
	b, err := json.Marshal(k)
	if err != nil {
		return ""
	}
	return string(b)
}

func (k *InlineKeyboard) AddButton(indexRow int, text string, callbackData string) {
	row := k.InlineKeyboardButton[indexRow]
	if row == nil {
		row[0] = InlineKeyboardButton{
			Text:         text,
			CallbackData: callbackData,
		}
	} else {
		row[len(row)] = InlineKeyboardButton{
			Text:         text,
			CallbackData: callbackData,
		}
	}
}

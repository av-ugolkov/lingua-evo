package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type CallbackQuery struct {
	Text    string           `json:"text"`
	Data    string           `json:"data"`
	From    From             `json:"from"`
	Chat    Chat             `json:"chat"`
	Message *IncomingMessage `json:"message"`
}

type From struct {
	ID       int    `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Username string `json:"username"`
}
type Chat struct {
	ID int `json:"id"`
}

package commands

const (
	UnknownCmd Command = "/unknown"
	StartCmd   Command = "/start"
	HelpCmd    Command = "/help"
	Cancel     Command = "/cancel"
	AddCmd     Command = "/add"
	RndCmd     Command = "/rnd"
)

type Command string

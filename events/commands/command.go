package commands

const (
	UnknownCmd = "/unknown"
	StartCmd   = "/start"
	HelpCmd    = "/help"
	Cancel     = "/cancel"
	AddCmd     = "/add"
	RndCmd     = "/rnd"
)

type Command string

type CommandExec interface {
	Execute(int) (Command, error)
}

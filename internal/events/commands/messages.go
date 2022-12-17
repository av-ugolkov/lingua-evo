package commands

const MsgHelp = `I can save and keep you pages. Also I can offer you then to read. 
In order to save the page, just send ne at link to it.
In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`
const MsgHello = "Hi %s! \n\n" + MsgHelp
const MsgCancel = "You cancel last command"

const (
	MsgUnknownCommand = "Unknown command"
	MsgNoSavedPages   = "You have no saved pages"
	MsgSaved          = "Saved!"
	MsgAlreadyExists  = "You have already heave this page in your list"
)

const (
	MsgAddWord      = "You need to write the word you want to learn"
	MsgAddPronounce = "How is the word pronounce?"
	MsgAddTranslate = "How is the word translated?"
	MsgAddExample   = "Example"
	MsgAddFinish    = "New word was added to dictionary"
)

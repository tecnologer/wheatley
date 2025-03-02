package commands

type CommandName string

const (
	Start                 CommandName = "start"
	AddStreamerCmdName    CommandName = "add"
	RemoveStreamerCmdName CommandName = "remove"
	ListStreamersCmdName  CommandName = "list"
	HelpCmdName           CommandName = "help"
)

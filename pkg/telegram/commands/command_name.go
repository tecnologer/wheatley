package commands

type CommandName string

const (
	Start          CommandName = "start"
	AddStreamer    CommandName = "add"
	RemoveStreamer CommandName = "remove"
)

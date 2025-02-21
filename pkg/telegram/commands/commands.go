package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao/db"
)

type (
	CommandHandler func(tgbotapi.Update, ...string) error
	HelpHandler    func() string
)

type Commands struct {
	Map map[CommandName]*Command
}

func NewCommands(db *db.Connection) *Commands {
	commands := &Commands{
		Map: make(map[CommandName]*Command),
	}

	commands.Add(
		AddStreamerCmd(db),
	)

	return commands
}

func (c *Commands) Add(cmds ...*Command) {
	for _, cmd := range cmds {
		c.Map[cmd.Name] = cmd
	}
}

func (c *Commands) IsRegistered(cmdName CommandName) bool {
	cmd, ok := c.Map[cmdName]

	return ok && cmd != nil
}

func (c *Commands) HasHandler(cmdName CommandName) bool {
	if !c.IsRegistered(cmdName) {
		return false
	}

	return c.Map[cmdName].Handler != nil
}

func (c *Commands) HasHelp(cmdName CommandName) bool {
	if !c.IsRegistered(cmdName) {
		return false
	}

	return c.Map[cmdName].Help != nil
}

func (c *Commands) Execute(cmdName string, update tgbotapi.Update, args ...string) error {
	cmd := CommandName(cmdName)

	if !c.HasHandler(cmd) {
		return nil
	}

	return c.Map[cmd].Handler(update, args...)
}

func (c *Commands) Help(cmdName string) string {
	cmd := CommandName(cmdName)

	if !c.HasHelp(cmd) {
		return ""
	}

	return c.Map[cmd].Help()
}

type Command struct {
	Name        CommandName
	Description string
	Handler     CommandHandler
	Help        HelpHandler
}

func DefaultCommands() *Commands {
	return &Commands{}
}

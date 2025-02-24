package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

var ErrHelpCmdNotFound = fmt.Errorf("help command not found")

func buildCmdHelpMessage(commands *Commands, cmdName CommandName) string {
	if !commands.IsRegistered(cmdName) {
		return ""
	}

	var (
		helpMsg strings.Builder
		command = commands.Map[cmdName]
	)

	helpMsg.WriteString("/")
	helpMsg.WriteString(string(command.Name))
	helpMsg.WriteString(": ")
	helpMsg.WriteString(command.Description)
	helpMsg.WriteString("\n\n")
	helpMsg.WriteString(command.Help())

	return helpMsg.String()
}

func helpMessageForCmdFromArgs(commands *Commands, argsOrder []string, args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}

	argsMapped, err := message.ArgsToMap(args, argsOrder)
	if err != nil {
		log.Errorf("parsing arguments: %v", err)

		return "", fmt.Errorf("the arguments are not valid")
	}

	if argsMapped["cmdName"] == "" {
		return "", nil
	}

	msg := buildCmdHelpMessage(commands, CommandName(argsMapped["cmdName"]))
	if msg == "" {
		return fmt.Sprintf("Command %s not found. Available commands:\n", argsMapped["cmdName"]), ErrHelpCmdNotFound
	}

	return msg, nil
}

func listCommandsSorted(commands *Commands) []CommandName {
	cmds := make([]CommandName, 0, len(commands.Map))

	for cmdName := range commands.Map {
		cmds = append(cmds, cmdName)
	}

	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i] < cmds[j]
	})

	return cmds
}

func MakeMarkdownLinkUser(username string) string {
	return fmt.Sprintf("[%s](https://twitch.tv/%s)", username, username)
}

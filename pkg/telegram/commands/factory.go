package commands

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

const (
	MsgStreamerAdded   = "Done! I'll notify you when `%s` goes live"
	MsgStreamerRemoved = "Done! You won't receive notifications for `%s` anymore"
)

var (
	ErrMissingStreamerName = fmt.Errorf("missing streamer name")
	ErrMissingChatID       = fmt.Errorf("missing chat ID")
)

func StartCmd() *Command {
	return &Command{
		Name:        Start,
		Description: "Starts the bot",
		Handler: func(cmd *Command, _ tgbotapi.Update, _ ...string) *Response {
			return NewResponse(WithCommand(cmd)).SetMessage("Hello! I'm Wheatley, your Twitch notifications bot")
		},
	}
}

func AddStreamerCmd(db *db.Connection) *Command {
	var (
		argsOrder = []string{"name"}
		daoNotif  = dao.NewNotifications(db)
	)

	return &Command{
		Name:        AddStreamer,
		Description: "Add a streamer to the list of notifications. You will be notified when the streamer goes live.",
		Help: func() string {
			return fmt.Sprintf("Usage: /%s <streamer_name", AddStreamer)
		},
		Handler: func(cmd *Command, update tgbotapi.Update, args ...string) *Response {
			response := NewResponse(WithCommand(cmd))

			if len(args) == 0 || args[0] == "" {
				log.Errorf("getting streamer name: %v", ErrMissingStreamerName)

				return response.SetMissingArgs("Missing streamer name")
			}

			chatID := message.GetChatIDFromUpdate(update)
			if chatID == 0 {
				log.Errorf("getting chat ID from update: %v", ErrMissingChatID)

				return response.SetMissingArgs("Missing chat ID")
			}

			argsMapped, err := message.ArgsToMap(args, argsOrder)
			if err != nil {
				log.Errorf("parsing arguments: %v", err)

				return response.SetMissingArgs("the arguments are not valid")
			}

			notification := &models.Notification{
				TwitchStreamerName: argsMapped["name"],
				TelegramChatID:     chatID,
			}

			err = daoNotif.CreateNotification(notification)
			if err != nil {
				log.Errorf("creating notification for streamer `%s`: %v", notification.TwitchStreamerName, err)

				return response.SetError("Error adding streamer %s", notification.TwitchStreamerName)
			}

			return response.SetMessage(MsgStreamerAdded, notification.TwitchStreamerName)
		},
	}
}

func RemoveStreamerCmd(db *db.Connection) *Command {
	var (
		argsOrder = []string{"name"}
		daoNotif  = dao.NewNotifications(db)
	)

	return &Command{
		Name:        RemoveStreamer,
		Description: "Remove a streamer from the list. You won't receive notifications for this streamer anymore.",
		Help: func() string {
			return fmt.Sprintf("Usage: /%s <streamer_name>", RemoveStreamer)
		},
		Handler: func(cmd *Command, update tgbotapi.Update, args ...string) *Response {
			response := NewResponse(WithCommand(cmd))

			if len(args) == 0 || args[0] == "" {
				log.Errorf("getting streamer name: %v", ErrMissingStreamerName)

				return response.SetMissingArgs("Missing streamer name")
			}

			chatID := message.GetChatIDFromUpdate(update)
			if chatID == 0 {
				log.Errorf("getting chat ID from update: %v", ErrMissingChatID)

				return response.SetMissingArgs("Missing chat ID")
			}

			argsMapped, err := message.ArgsToMap(args, argsOrder)
			if err != nil {
				log.Errorf("parsing arguments: %v", err)

				return response.SetMissingArgs("the arguments are not valid")
			}

			notification := &models.Notification{
				TwitchStreamerName: argsMapped["name"],
				TelegramChatID:     chatID,
			}

			err = daoNotif.DeleteNotification(notification)
			if err != nil {
				log.Errorf("deleting notification for streamer `%s`: %v", notification.TwitchStreamerName, err)

				return response.SetError("Error removing streamer %s", notification.TwitchStreamerName)
			}

			return response.SetMessage(MsgStreamerRemoved, notification.TwitchStreamerName)
		},
	}
}

func HelpCmd(commands *Commands) *Command {
	argsOrder := []string{"cmdName"}

	return &Command{
		Name:        Help,
		Description: "Shows the available commands and their usage.",
		Help: func() string {
			return fmt.Sprintf("Usage: /%s <command>", Help)
		},
		Handler: func(cmd *Command, _ tgbotapi.Update, args ...string) *Response {
			response := NewResponse(WithCommand(cmd))

			var helpMsg strings.Builder

			helpCmdMsg, err := helpMessageForCmdFromArgs(commands, argsOrder, args)
			if helpCmdMsg != "" && err == nil {
				return response.SetMessage(helpCmdMsg)
			}

			if helpCmdMsg != "" {
				helpMsg.WriteString(helpCmdMsg)
				helpMsg.WriteString("\n")
			}

			if err != nil && !errors.Is(err, ErrHelpCmdNotFound) {
				helpMsg.WriteString(err.Error())
				helpMsg.WriteString("\n")
			}

			helpMsg.WriteString("Hello! I'm Wheatley, your Twitch notifications bot, and I can help you with the following commands:\n\n")

			for _, cmdName := range listCommandsSorted(commands) {
				cmd := commands.Map[cmdName]

				helpMsg.WriteString("» /")
				helpMsg.WriteString(string(cmd.Name))
				helpMsg.WriteString(": ")
				helpMsg.WriteString(cmd.Description)
				helpMsg.WriteString("\n")
			}

			helpMsg.WriteString("\n")
			helpMsg.WriteString("For more information about a specific command, use /help <command>. For example: /help add")

			return response.SetMessage(helpMsg.String())
		},
	}
}

func ListStreamersCmd(db *db.Connection) *Command {
	daoNotif := dao.NewNotifications(db)

	return &Command{
		Name:        ListStreamers,
		Description: "List the streamers you are subscribed to",
		Handler: func(cmd *Command, update tgbotapi.Update, _ ...string) *Response {
			response := NewResponse(WithCommand(cmd))

			chatID := message.GetChatIDFromUpdate(update)
			if chatID == 0 {
				log.Errorf("getting chat ID from update: %v", ErrMissingChatID)

				return response.SetMissingArgs("Missing chat ID")
			}

			notifications, err := daoNotif.NotificationsByChatID(chatID)
			if err != nil {
				log.Errorf("getting notifications: %v", err)

				return response.SetError("Error getting streamers")
			}

			if len(notifications) == 0 {
				return response.SetMessage("You are not subscribed to any streamer yet")
			}

			var msg strings.Builder

			msg.WriteString("You are subscribed to the following streamers:\n\n")

			for _, notif := range notifications {
				msg.WriteString("🎮 ")
				msg.WriteString(MakeMarkdownLinkUser(notif.TwitchStreamerName))
				msg.WriteString("\n")
			}

			return response.SetMessage(msg.String())
		},
	}
}

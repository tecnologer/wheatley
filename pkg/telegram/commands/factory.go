package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

const (
	MsgStreamerAdded   = "Done! I'll notify you when %s goes live"
	MsgStreamerRemoved = "Done! You won't receive notifications for %s anymore"
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
		Description: "Add a streamer to the list",
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
				log.Errorf("creating notification for streamer %s: %v", notification.TwitchStreamerName, err)

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
		Description: "Remove a streamer from the list",
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
				log.Errorf("deleting notification for streamer %s: %v", notification.TwitchStreamerName, err)

				return response.SetError("Error removing streamer %s", notification.TwitchStreamerName)
			}

			return response.SetMessage(MsgStreamerRemoved, notification.TwitchStreamerName)
		},
	}
}

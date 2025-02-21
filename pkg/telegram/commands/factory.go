package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

var (
	ErrMissingArgs         = fmt.Errorf("missing arguments")
	ErrMissingStreamerName = fmt.Errorf("missing streamer name")
	ErrMissingChatID       = fmt.Errorf("missing chat ID")
)

func AddStreamerCmd(db *db.Connection) *Command {
	var (
		argsOrder = []string{"name"}
		daoNotif  = dao.NewNotifications(db)
	)

	return &Command{
		Name:        AddStreamer,
		Description: "Add a streamer to the list",
		Handler: func(update tgbotapi.Update, args ...string) error {
			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("%w: %w", ErrMissingArgs, ErrMissingStreamerName)
			}

			chatID := message.GetChatIDFromUpdate(update)
			if chatID == 0 {
				return fmt.Errorf("%w: %w", ErrMissingArgs, ErrMissingChatID)
			}

			argsMapped, err := message.ArgsToMap(args, argsOrder)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrMissingArgs, err)
			}

			notification := &models.Notification{
				TwitchStreamerName: argsMapped["name"],
				TelegramChatID:     chatID,
			}

			err = daoNotif.CreateNotification(notification)
			if err != nil {
				return fmt.Errorf("creating notification for streamer %s: %w", notification.TwitchStreamerName, err)
			}

			return nil
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
			return "Usage: /remove <streamer_name>"
		},
		Handler: func(update tgbotapi.Update, args ...string) error {
			if len(args) == 0 || args[0] == "" {
				return fmt.Errorf("%w: %w", ErrMissingArgs, ErrMissingStreamerName)
			}

			chatID := message.GetChatIDFromUpdate(update)
			if chatID == 0 {
				return fmt.Errorf("%w: %w", ErrMissingArgs, ErrMissingChatID)
			}

			argsMapped, err := message.ArgsToMap(args, argsOrder)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrMissingArgs, err)
			}

			notification := &models.Notification{
				TwitchStreamerName: argsMapped["name"],
				TelegramChatID:     chatID,
			}

			err = daoNotif.DeleteNotification(notification)
			if err != nil {
				return fmt.Errorf("deleting notification for streamer %s: %w", notification.TwitchStreamerName, err)
			}

			return nil
		},
	}
}

package commands

import (
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/tecnologer/wheatley/pkg/constants"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

// notifyAdminAddedStreamer creates the response to notify the admin when a user adds a streamer to the list
func notifyAdminAddedStreamer(cmd *Command, update tgbotapi.Update, args ...string) *Response {
	var msg strings.Builder

	response := NewResponse(WithCommand(cmd))

	msg.WriteString("The ")

	if message.IsBot(update) {
		msg.WriteString("bot @")
	} else {
		msg.WriteString("user @")
	}

	msg.WriteString(message.SenderName(update))

	msg.WriteString(" added the streamer `")
	msg.WriteString(args[0])
	msg.WriteString("` to chat \"")
	msg.WriteString(message.GetChatNameFromUpdate(update))
	msg.WriteString("\".")

	return response.SetMessage(msg.String())
}

// notifyAdminRemovedStreamer creates the response to notify the admin when a user removes a streamer from the list
func notifyAdminRemovedStreamer(cmd *Command, update tgbotapi.Update, args ...string) *Response {
	var msg strings.Builder

	response := NewResponse(WithCommand(cmd))

	msg.WriteString("The ")

	if message.IsBot(update) {
		msg.WriteString("bot @")
	} else {
		msg.WriteString("user @")
	}

	msg.WriteString(message.SenderName(update))

	msg.WriteString(" removed the streamer `")
	msg.WriteString(args[0])
	msg.WriteString("` from chat \"")
	msg.WriteString(message.GetChatNameFromUpdate(update))
	msg.WriteString("\".")

	return response.SetMessage(msg.String())
}

func buildNotificationFromUpdate(update tgbotapi.Update, args, argsOrder []string) (*models.Notification, error) {
	if len(args) == 0 || args[0] == "" {
		return nil, fmt.Errorf("missing streamer name")
	}

	chatID := message.GetChatIDFromUpdate(update)
	if chatID == 0 {
		return nil, fmt.Errorf("missing chat ID")
	}

	argsMapped, err := message.ArgsToMap(args, argsOrder)
	if err != nil {
		return nil, fmt.Errorf("the arguments are not valid. %w", err)
	}

	if argsMapped["name"] == "" {
		return nil, fmt.Errorf("missing streamer name in arguments")
	}

	var threadID *int
	if updateThreadID := message.GetMessageThreadID(update); updateThreadID != 0 {
		threadID = &updateThreadID
	}

	argsMapped["name"] = strings.TrimPrefix(strings.ToLower(argsMapped["name"]), constants.TwitchURLPrefix)

	notification := &models.Notification{
		TwitchStreamerName: argsMapped["name"],
		TelegramChatID:     chatID,
		TelegramThreadID:   threadID,
	}

	return notification, nil
}

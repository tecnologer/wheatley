package message

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetFromUpdate(update tgbotapi.Update) string {
	if update.Message == nil && update.EditedMessage == nil {
		return ""
	}

	if update.Message == nil && update.EditedMessage != nil {
		return update.EditedMessage.Text
	}

	return update.Message.Text
}

func GetChatIDFromUpdate(update tgbotapi.Update) int64 {
	if update.Message == nil && update.EditedMessage == nil {
		return 0
	}

	if update.Message != nil && update.Message.Chat != nil {
		return update.Message.Chat.ID
	}

	if update.EditedMessage != nil && update.EditedMessage.Chat != nil {
		return update.EditedMessage.Chat.ID
	}

	return 0
}

func GetChatNameFromUpdate(update tgbotapi.Update) string {
	if update.Message == nil && update.EditedMessage == nil {
		return "<<no defined>>"
	}

	if update.Message != nil && update.Message.Chat != nil {
		if update.Message.Chat.UserName != "" {
			return update.Message.Chat.UserName
		}

		return update.Message.Chat.Title
	}

	if update.EditedMessage != nil && update.EditedMessage.Chat != nil {
		if update.EditedMessage.Chat.UserName != "" {
			return update.EditedMessage.Chat.UserName
		}

		return update.EditedMessage.Chat.Title
	}

	return "<<no defined>>"
}

func IsBot(update tgbotapi.Update) bool {
	if update.Message == nil && update.EditedMessage == nil {
		return false
	}

	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.IsBot
	}

	if update.EditedMessage != nil && update.EditedMessage.From != nil {
		return update.EditedMessage.From.IsBot
	}

	return false
}

func SenderName(update tgbotapi.Update) string {
	if update.Message == nil && update.EditedMessage == nil {
		return ""
	}

	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.UserName
	}

	if update.EditedMessage != nil && update.EditedMessage.From != nil {
		return update.EditedMessage.From.UserName
	}

	return ""
}

// SentByAdmin returns true if the message was sent by an admin
func SentByAdmin(update tgbotapi.Update, admins []string) bool {
	author := SenderName(update)

	for _, admin := range admins {
		if strings.EqualFold(author, admin) {
			return true
		}
	}

	return false
}

package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

	if update.Message == nil && update.EditedMessage != nil {
		return update.EditedMessage.Chat.ID
	}

	return update.Message.Chat.ID
}

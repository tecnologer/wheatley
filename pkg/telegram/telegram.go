package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

type Bot struct {
	*tgbotapi.BotAPI
	commands *commands.Commands
}

func NewBot(token string, verbose bool, dbCnn *db.Connection) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	bot.Debug = verbose

	log.Infof("authorized on account %s", bot.Self.UserName)

	return &Bot{
		BotAPI:   bot,
		commands: commands.NewCommands(dbCnn),
	}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := b.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (b *Bot) ReadUpdates() error {
	var err error

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		err = b.ExecCommand(update)
		if err != nil {
			b.SendMessage(update.Message.Chat.ID, fmt.Sprintf("error: %s", err))
		}
	}

	return nil
}

func (b *Bot) ExecCommand(update tgbotapi.Update) error {
	msg := message.GetFromUpdate(update)
	if msg == "" {
		return nil
	}

	cmdName, args := message.ExtractCommand(msg)
	if cmdName == "" {
		return nil
	}

	err := b.commands.Execute(cmdName, update, args...)
	if err == nil {
		b.SendMessage(update.Message.Chat.ID, "command executed")
		return nil
	}

	return nil
}

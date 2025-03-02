package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

type Config struct {
	Token     string
	Verbose   bool
	DB        *db.Connection
	ChatAdmin int64
	Admins    []string
}

func (c *Config) String() string {
	if c == nil {
		return "<<no config>>"
	}

	var token string
	if len(c.Token) > 5 {
		token = c.Token[:5]
	}

	return fmt.Sprintf("Token: ...%s, Verbose: %t, ChatAdmin: %d, Admins: %v", token, c.Verbose, c.ChatAdmin, c.Admins)
}

type Bot struct {
	*Config
	*tgbotapi.BotAPI
	commands *commands.Commands
}

func NewBot(config *Config) (*Bot, error) {
	log.Infof("creating bot with config: %s", config)

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	bot.Debug = config.Verbose

	log.Infof("authorized on account %s", bot.Self.UserName)

	return &Bot{
		Config:   config,
		BotAPI:   bot,
		commands: commands.NewCommands(config.DB),
	}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := b.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (b *Bot) ReadUpdates() error {
	var (
		err error
		msg string
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		msg = b.ExecCommand(update)
		if msg == "" {
			continue
		}

		err = b.SendMessage(message.GetChatIDFromUpdate(update), msg)
		if err != nil {
			log.Errorf("sending message: %v", err)
		}

		go b.NotifyAdminIfNecessary(update)
	}

	return nil
}

func (b *Bot) ExecCommand(update tgbotapi.Update) string {
	inputMsg := message.GetFromUpdate(update)
	if inputMsg == "" {
		return ""
	}

	cmdName, args := b.extractCommand(inputMsg)
	if cmdName == "" {
		return ""
	}

	return b.commands.Execute(cmdName, update, args...)
}

func (b *Bot) extractCommand(inputMsg string) (string, []string) {
	cmdName, args := message.ExtractCommand(inputMsg)
	if cmdName == "" {
		return "", nil
	}

	return strings.ReplaceAll(cmdName, "@"+b.Self.UserName, ""), args
}

func (b *Bot) NotifyAdminIfNecessary(update tgbotapi.Update) {
	cmdName, args := b.extractCommand(message.GetFromUpdate(update))
	if cmdName == "" {
		return
	}

	if !b.shouldNotifyAdminChat(update, cmdName) {
		return
	}

	res := b.commands.AdminNotification(cmdName, update, args...)

	err := b.SendMessage(b.ChatAdmin, res.Message())
	if err != nil {
		log.Errorf("sending message to admin: %v", err)
	}
}

func (b *Bot) shouldNotifyAdminChat(update tgbotapi.Update, cmdName string) bool {
	return b.ChatAdmin != 0 && // if the admin chat is defined
		message.GetChatIDFromUpdate(update) != b.ChatAdmin && // if the chat is not the admin chat
		!message.SentByAdmin(update, b.Admins) && // if the author message is not an admin
		b.commands.HasAdminNotification(cmdName) // if the command has an admin notification
}

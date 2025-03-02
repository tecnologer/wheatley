package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Self() tgbotapi.User
}

type BotAPIImpl struct {
	*tgbotapi.BotAPI
}

func (b *BotAPIImpl) Self() tgbotapi.User {
	return b.BotAPI.Self
}

type Config struct {
	Token     string
	Verbose   bool
	DB        *db.Connection
	ChatAdmin int64
	Admins    []string
	IsMock    bool
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

func (c *Config) OK() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.Token == "" {
		return fmt.Errorf("missing token")
	}

	if !c.IsMock && c.DB == nil {
		return fmt.Errorf("missing database connection")
	}

	return nil
}

type Bot struct {
	*Config
	Bot      BotAPI
	commands *commands.Commands
}

func NewBot(config *Config) (*Bot, error) {
	if err := config.OK(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	log.Infof("creating bot with config: %s", config)

	bot := BotAPI(&Mock{})

	if !config.IsMock {
		tBot, err := tgbotapi.NewBotAPI(config.Token)
		if err != nil {
			return nil, fmt.Errorf("creating bot: %w", err)
		}

		tBot.Debug = config.Verbose

		bot = &BotAPIImpl{tBot}
	}

	log.Infof("authorized on account %s", bot.Self().UserName)

	return &Bot{
		Config:   config,
		Bot:      bot,
		commands: commands.NewCommands(config.DB),
	}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := b.Bot.Send(msg)
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

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		log.Info("received update")

		msg = b.ExecCommand(update)
		if msg == "" {
			continue
		}

		log.Infof("received message: %s", msg)

		err = b.SendMessage(message.GetChatIDFromUpdate(update), msg)
		if err != nil {
			log.Errorf("sending message: %v", err)
		}

		log.Info("message sent")

		go b.NotifyAdminIfNecessary(update)
	}

	return nil
}

func (b *Bot) ExecCommand(update tgbotapi.Update) string {
	inputMsg := message.GetFromUpdate(update)
	if inputMsg == "" {
		return ""
	}

	cmdName, args := message.ExtractCommandNamedBot(inputMsg, b.Bot.Self().UserName)
	if cmdName == "" {
		return ""
	}

	return b.commands.Execute(cmdName, update, args...)
}

func (b *Bot) NotifyAdminIfNecessary(update tgbotapi.Update) {
	log.Infof("notifyAdminIfNecessary update")

	cmdName, args := message.ExtractCommandNamedBot(message.GetFromUpdate(update), b.Bot.Self().UserName)
	if cmdName == "" {
		return
	}

	log.Infof("notifyAdminIfNecessary command: %s", cmdName)

	if !b.shouldNotifyAdminChat(update, cmdName) {
		return
	}

	log.Infof("notifyAdminIfNecessary should notify admin")

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

package cli

import (
	"context"
	"fmt"

	"github.com/tecnologer/wheatley/cmd/flags"
	"github.com/tecnologer/wheatley/pkg/cron"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/telegram"
	"github.com/tecnologer/wheatley/pkg/twitch"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"github.com/urfave/cli/v2"
)

type CLI struct {
	*cli.App
	bot *telegram.Bot
}

func NewCLI(versionValue string) *CLI {
	newCLI := &CLI{}

	newCLI.setupApp(versionValue)

	return newCLI
}

func (c *CLI) setupApp(versionValue string) {
	c.App = &cli.App{
		Name:        "wheatley",
		Version:     versionValue,
		Usage:       "Execute the bot to interact with the Twitch API and the users to send notifications when a streamer goes live.",
		Description: "",
		Action:      c.run,
		Before:      c.beforeRun,
		Flags: []cli.Flag{
			flags.TelegramToken(),
			flags.Interval(),
			flags.Verbose(),
			flags.DBName(),
			flags.DBPassword(),
			flags.TwitchClientID(),
			flags.TwitchClientSecret(),
			flags.ResendInterval(),
			flags.TelegramAdminChatID(),
			flags.TelegramAdmins(),
		},
		EnableBashCompletion: true,
	}
}

func (c *CLI) beforeRun(ctx *cli.Context) error {
	// Disable color globally.
	if ctx.Bool(flags.VerboseFlagName) {
		log.SetLevel(log.DebugLevel)
	}

	if ctx.String(flags.TelegramTokenFlagName) == "" {
		return fmt.Errorf("telegram token is required")
	}

	if ctx.String(flags.DBNameFlagName) == "" {
		return fmt.Errorf("db name is required")
	}

	return nil
}

func (c *CLI) run(ctx *cli.Context) error {
	log.Info("creating context")

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Infof("creating db connection")

	dbCnn, err := c.createConnection(ctx)
	if err != nil {
		return fmt.Errorf("create db connection: %w", err)
	}

	log.Infof("creating telegram bot")

	c.bot, err = telegram.NewBot(&telegram.Config{
		Token:     ctx.String(flags.TelegramTokenFlagName),
		Verbose:   ctx.Bool(flags.VerboseFlagName),
		DB:        dbCnn,
		ChatAdmin: ctx.Int64(flags.TelegramAdminChatIDsFlagName),
		Admins:    ctx.StringSlice(flags.TelegramAdminsFlagName),
	})
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	log.Infof("creating cron job")

	cronJob, err := cron.NewScheduler(c.cronOptions(runCtx, ctx, dbCnn))
	if err != nil {
		return fmt.Errorf("create cron job: %w", err)
	}

	cronJob.Start()

	defer func() {
		err := cronJob.Shutdown()
		if err != nil {
			log.Errorf("shutdown cron job: %v", err)
		}
	}()

	log.Infof("cron job started")

	log.Infof("telegram bot started, reading updates")

	err = c.bot.ReadUpdates()
	if err != nil {
		return fmt.Errorf("read updates: %w", err)
	}

	return nil
}

func (c *CLI) createConnection(ctx *cli.Context) (*db.Connection, error) {
	log.Infof("connecting to the DB at %s", ctx.String(flags.DBNameFlagName))

	dbConfig := &db.Config{
		Password: ctx.String(flags.DBPasswordFlagName),
		DBName:   ctx.String(flags.DBNameFlagName),
	}

	cnn, err := db.NewConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("create new connection: %w", err)
	}

	log.Infof("connection to the DB at %s established", ctx.String(flags.DBNameFlagName))

	return cnn, nil
}

func (c *CLI) cronOptions(runCtx context.Context, ctx *cli.Context, dbCnn *db.Connection) *cron.Config {
	return &cron.Config{
		IntervalMinutes:        ctx.Int(flags.IntervalFlagName),
		NotificationDelayHours: ctx.Int(flags.ResendIntervalFlagName),
		TwitchConfig: &twitch.Config{
			ClientID:     ctx.String(flags.TwitchClientIDFlagName),
			ClientSecret: ctx.String(flags.TwitchClientSecretFlagName),
		},
		Notifications: dao.NewNotifications(dbCnn),
		Context:       runCtx,
		TelegramBot:   c.bot,
	}
}

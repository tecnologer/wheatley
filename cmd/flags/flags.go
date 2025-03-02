package flags

import (
	"github.com/tecnologer/wheatley/pkg/contants/envvarname"
	"github.com/urfave/cli/v2"
)

const (
	VerboseFlagName              = "verbose"
	DBPasswordFlagName           = "db-password"
	DBNameFlagName               = "db-name"
	TelegramTokenFlagName        = "telegram-token"
	IntervalFlagName             = "interval"
	ResendIntervalFlagName       = "resend-interval"
	TwitchClientIDFlagName       = "twitch-client-id"
	TwitchClientSecretFlagName   = "twitch-client-secret"
	TelegramAdminChatIDsFlagName = "telegram-admin-chat-ids"
	TelegramAdminsFlagName       = "telegram-admins"
)

func Verbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    VerboseFlagName,
		Usage:   "Enable verbose output.",
		Aliases: []string{"V"},
	}
}

func DBPassword() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBPasswordFlagName,
		Aliases:  []string{"p"},
		Usage:    "Database password.",
		Required: false,
	}
}

func DBName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBNameFlagName,
		Aliases:  []string{"d"},
		Usage:    "Database name.",
		Required: true,
		EnvVars:  []string{envvarname.DBName},
	}
}

func TelegramToken() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     TelegramTokenFlagName,
		Aliases:  []string{"t"},
		Usage:    "Telegram bot token.",
		Required: true,
		EnvVars:  []string{envvarname.TelegramBotToken},
	}
}

func Interval() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    IntervalFlagName,
		Aliases: []string{"i"},
		Usage:   "Interval in minutes to check if a streamer is live.",
		Value:   1,
		EnvVars: []string{envvarname.Interval},
	}
}

func TwitchClientID() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    TwitchClientIDFlagName,
		Aliases: []string{"c"},
		Usage:   "Twitch client ID.",
		EnvVars: []string{envvarname.TwitchClientID},
	}
}

func TwitchClientSecret() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    TwitchClientSecretFlagName,
		Aliases: []string{"s"},
		Usage:   "Twitch client secret.",
		EnvVars: []string{envvarname.TwitchClientSecret},
	}
}

func ResendInterval() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    ResendIntervalFlagName,
		Aliases: []string{"r"},
		Usage:   "Interval in hours to resend a notification.",
		Value:   6,
		EnvVars: []string{envvarname.ResendInterval},
	}
}

func TelegramAdminChatID() *cli.Int64Flag {
	return &cli.Int64Flag{
		Name:    TelegramAdminChatIDsFlagName,
		Usage:   "The main Telegram chat ID of the admin.",
		Aliases: []string{"D"},
		Value:   10244644, // Default chat ID for the bot and the owner
		EnvVars: []string{envvarname.TelegramAdminChatIDs},
	}
}

func TelegramAdmins() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    TelegramAdminsFlagName,
		Usage:   "The telegram username of the admins.",
		Aliases: []string{"a"},
		Value:   cli.NewStringSlice("tecnologer"),
		EnvVars: []string{envvarname.TelegramAdminChatIDs},
	}
}

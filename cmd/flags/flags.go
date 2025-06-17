package flags

import (
	"github.com/tecnologer/wheatley/pkg/constants/envvarname"
	"github.com/urfave/cli/v2"
)

const (
	VerboseFlagName               = "verbose"
	VerboseFlagAlias              = "V"
	DBPasswordFlagName            = "db-password"
	DBPasswordFlagAlias           = "p"
	DBNameFlagName                = "db-name"
	DBNameFlagAlias               = "d"
	TelegramTokenFlagName         = "telegram-token"
	TelegramTokenFlagAlias        = "t"
	IntervalFlagName              = "interval"
	IntervalFlagAlias             = "i"
	ResendIntervalFlagName        = "resend-interval"
	ResendIntervalFlagAlias       = "r"
	TwitchClientIDFlagName        = "twitch-client-id"
	TwitchClientIDFlagAlias       = "c"
	TwitchClientSecretFlagName    = "twitch-client-secret"
	TwitchClientSecretFlagAlias   = "s"
	TelegramAdminChatIDsFlagName  = "telegram-admin-chat-ids"
	TelegramAdminChatIDsFlagAlias = "D"
	TelegramAdminsFlagName        = "telegram-admins"
	TelegramAdminsFlagAlias       = "a"
)

func Verbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    VerboseFlagName,
		Usage:   "Enable verbose output.",
		Aliases: []string{VerboseFlagAlias},
	}
}

func DBPassword() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBPasswordFlagName,
		Aliases:  []string{DBPasswordFlagAlias},
		Usage:    "Database password.",
		Required: false,
	}
}

func DBName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBNameFlagName,
		Aliases:  []string{DBNameFlagAlias},
		Usage:    "Database name.",
		Required: true,
		EnvVars:  []string{envvarname.DBName},
	}
}

func TelegramToken() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     TelegramTokenFlagName,
		Aliases:  []string{TelegramTokenFlagAlias},
		Usage:    "Telegram bot token.",
		Required: true,
		EnvVars:  []string{envvarname.TelegramBotToken},
	}
}

func Interval() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    IntervalFlagName,
		Aliases: []string{IntervalFlagAlias},
		Usage:   "Interval in minutes to check if a streamer is live.",
		Value:   1,
		EnvVars: []string{envvarname.Interval},
	}
}

func TwitchClientID() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    TwitchClientIDFlagName,
		Aliases: []string{TwitchClientIDFlagAlias},
		Usage:   "Twitch client ID.",
		EnvVars: []string{envvarname.TwitchClientID},
	}
}

func TwitchClientSecret() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    TwitchClientSecretFlagName,
		Aliases: []string{TwitchClientSecretFlagAlias},
		Usage:   "Twitch client secret.",
		EnvVars: []string{envvarname.TwitchClientSecret},
	}
}

func ResendInterval() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    ResendIntervalFlagName,
		Aliases: []string{ResendIntervalFlagAlias},
		Usage:   "Interval in hours to resend a notification.",
		Value:   6,
		EnvVars: []string{envvarname.ResendInterval},
	}
}

func TelegramAdminChatID() *cli.Int64Flag {
	return &cli.Int64Flag{
		Name:    TelegramAdminChatIDsFlagName,
		Usage:   "The main Telegram chat ID of the admin.",
		Aliases: []string{TelegramAdminChatIDsFlagAlias},
		Value:   10244644, // Default chat ID for the bot and the owner
		EnvVars: []string{envvarname.TelegramAdminChatIDs},
	}
}

func TelegramAdmins() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    TelegramAdminsFlagName,
		Usage:   "The telegram username of the admins.",
		Aliases: []string{TelegramAdminsFlagAlias},
		Value:   cli.NewStringSlice("tecnologer"),
		EnvVars: []string{envvarname.TelegramAdminChatIDs},
	}
}

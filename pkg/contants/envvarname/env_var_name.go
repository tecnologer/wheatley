package envvarname

const (
	DBPassword           = "WHEATLEY_DB_PASSWORD"
	DBName               = "WHEATLEY_DB_NAME"
	TelegramBotToken     = "WHEATLEY_TELEGRAM_BOT_TOKEN" //nolint:gosec // This is not a real token
	Interval             = "WHEATLEY_INTERVAL"
	ResendInterval       = "WHEATLEY_RESEND_INTERVAL"
	TwitchClientID       = "WHEATLEY_TWITCH_CLIENT_ID"
	TwitchClientSecret   = "WHEATLEY_TWITCH_CLIENT_SECRET"
	TelegramAdminChatIDs = "WHEATLEY_TELEGRAM_ADMIN_CHAT_IDS"
)

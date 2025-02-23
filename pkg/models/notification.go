package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	*gorm.Model
	TelegramChatID     int64     `json:"telegram_id"          gorm:"index:idx_telegram_chat_id_twitch_streamer_name,unique"`
	TwitchStreamerName string    `json:"twitch_streamer_name" gorm:"index:idx_telegram_chat_id_twitch_streamer_name,unique"`
	LastNotification   time.Time `json:"last_notification"    gorm:"default:null"`
}

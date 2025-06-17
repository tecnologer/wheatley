package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	*gorm.Model
	TelegramChatID     int64     `json:"telegram_id"          gorm:"index:idx_telegram_chat_id_twitch_streamer_name,unique,,priority:1"`
	TwitchStreamerName string    `json:"twitch_streamer_name" gorm:"index:idx_telegram_chat_id_twitch_streamer_name,unique,,priority:2"`
	LastNotification   time.Time `json:"last_notification"    gorm:"default:null"`
	LastGame           string    `json:"last_game"            gorm:"default:null"`
	TelegramThreadID   *int      `json:"telegram_thread_id"   gorm:"index:idx_telegram_chat_id_twitch_streamer_name,unique,,priority:3"`
}

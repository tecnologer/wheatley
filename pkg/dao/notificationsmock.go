package dao

import (
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/models"
)

type NotificationsDAO interface {
	CreateNotification(notification *models.Notification) error
	DeleteNotification(notification *models.Notification) error
	NotificationsByChatID(chatID int64) ([]*models.Notification, error)
}

func NewNotificationsDAO(db *db.Connection) NotificationsDAO {
	if db == nil {
		return NotificationsMock{}
	}

	return NewNotifications(db)
}

type NotificationsMock struct{}

func (n NotificationsMock) CreateNotification(_ *models.Notification) error {
	return nil
}

func (n NotificationsMock) DeleteNotification(_ *models.Notification) error {
	return nil
}

func (n NotificationsMock) NotificationsByChatID(chatID int64) ([]*models.Notification, error) {
	return []*models.Notification{
		{
			TwitchStreamerName: "streamer_name",
			TelegramChatID:     chatID,
		},
		{
			TwitchStreamerName: "another_streamer",
			TelegramChatID:     chatID,
		},
	}, nil
}

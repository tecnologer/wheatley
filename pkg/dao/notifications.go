package dao

import (
	"errors"
	"fmt"

	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/models"
	"gorm.io/gorm"
)

type Notifications struct {
	db *db.Connection
}

func NewNotifications(db *db.Connection) *Notifications {
	return &Notifications{
		db: db,
	}
}

func (s *Notifications) NotificationByStreamerName(chatID int64, streamerName string) (*models.Notification, error) {
	var streamer models.Notification

	err := s.db.Where("twitch_streamer_name = ? AND telegram_chat_id = ?", streamerName, chatID).First(&streamer).Error
	if err != nil {
		return nil, fmt.Errorf("getting notification settings: %w", err)
	}

	return &streamer, nil
}

func (s *Notifications) CreateNotification(notification *models.Notification) error {
	existing, err := s.NotificationByStreamerName(notification.TelegramChatID, notification.TwitchStreamerName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("getting existing notification settings: %w", err)
	}

	// If the notification already exists, we don't need to create it again
	if existing != nil {
		return nil
	}

	err = s.db.Create(notification).Error
	if err != nil {
		return fmt.Errorf("creating notification settings: %w", err)
	}

	return nil
}

func (s *Notifications) DeleteNotification(notification *models.Notification) error {
	err := s.db.Unscoped().
		Where("twitch_streamer_name = ? AND telegram_chat_id = ?", notification.TwitchStreamerName, notification.TelegramChatID).
		Delete(notification).
		Error
	if err != nil {
		return fmt.Errorf("deleting notification settings: %w", err)
	}

	return nil
}

func (s *Notifications) AllNotifications() ([]*models.Notification, error) {
	var notifications []*models.Notification

	err := s.db.Find(&notifications).Error
	if err != nil {
		return nil, fmt.Errorf("getting all notifications: %w", err)
	}

	return notifications, nil
}

func (s *Notifications) UpdateNotification(notification *models.Notification) error {
	err := s.db.Save(notification).Error
	if err != nil {
		return fmt.Errorf("updating notification settings: %w", err)
	}

	return nil
}

func (s *Notifications) NotificationsByChatID(chatID int64) ([]*models.Notification, error) {
	var notifications []*models.Notification

	err := s.db.Where("telegram_chat_id = ?", chatID).Find(&notifications).Error
	if err != nil {
		return nil, fmt.Errorf("getting notifications by chat ID: %w", err)
	}

	return notifications, nil
}

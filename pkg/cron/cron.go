package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/go-co-op/gocron/v2"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/telegram"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
	"github.com/tecnologer/wheatley/pkg/twitch"
	"github.com/tecnologer/wheatley/pkg/utils/log"
)

type Config struct {
	IntervalMinutes        int
	NotificationDelayHours int
	TwitchConfig           *twitch.Config
	Notifications          dao.NotificationsDAO
	Context                context.Context //nolint: containedctx // This is a context.Context that will be used to make requests to Twitch API
	TelegramBot            *telegram.Bot
}

type Scheduler struct {
	gocron.Scheduler
	*Config
	twitch twitch.API
}

func NewScheduler(config *Config) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("creating scheduler: %w", err)
	}

	schedule := &Scheduler{
		Scheduler: scheduler,
		Config:    config,
		twitch:    twitch.New(config.TwitchConfig),
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Duration(config.IntervalMinutes)*time.Minute),
		gocron.NewTask(schedule.taskTwitchCheckStreamers),
	)
	if err != nil {
		return nil, fmt.Errorf("creating job: %w", err)
	}

	return schedule, nil
}

func (s *Scheduler) taskTwitchCheckStreamers() {
	log.Info("checking streamers")

	notifications, err := s.Notifications.AllNotifications()
	if err != nil {
		log.Errorf("getting notifications: %v", err)

		return
	}

	for _, notification := range notifications {
		stream, err := s.twitch.StreamByName(s.Context, notification.TwitchStreamerName)
		if err != nil {
			s.manageStreamerErr(err, notification)

			continue
		}

		if stream == nil || s.TelegramBot == nil || !s.requireSendMessage(notification, stream.GameName) {
			continue
		}

		notification.LastGame = stream.GameName

		s.sendMessage(stream, notification)
	}
}

func (s *Scheduler) manageStreamerErr(err error, notification *models.Notification) {
	if errors.Is(err, twitch.ErrNotFound) {
		if !notification.LastNotification.IsZero() {
			s.notifyStreamerWentOffline(notification)
		}
	}

	log.Errorf("getting stream for %s: %v", notification.TwitchStreamerName, err)
}

func (s *Scheduler) notifyStreamerWentOffline(notification *models.Notification) {
	err := s.TelegramBot.SendMessage(
		notification.TelegramChatID,
		fmt.Sprintf("Streamer `%s` went offline", notification.TwitchStreamerName),
	)
	if err != nil {
		log.Errorf("sending message: %v", err)
	}

	notification.LastNotification = time.Time{}
	s.updateNotification(notification)
}

func (s *Scheduler) updateNotification(notification *models.Notification) {
	err := s.Notifications.UpdateNotification(notification)
	if err != nil {
		log.Errorf("updating notification for %s: %v", notification.TwitchStreamerName, err)
	}
}

func (s *Scheduler) requireSendMessage(notification *models.Notification, currentGame string) bool {
	if notification.LastNotification.IsZero() {
		return true
	}

	return time.Since(notification.LastNotification) >= time.Duration(s.NotificationDelayHours)*time.Hour ||
		notification.LastGame != currentGame
}

func (s *Scheduler) sendMessage(stream *api.Stream, notification *models.Notification) {
	err := s.TelegramBot.SendMessage(
		notification.TelegramChatID,
		s.buildMessage(stream),
	)
	if err != nil {
		log.Errorf("sending message for online streamer %s: %v", notification.TwitchStreamerName, err)
	}

	notification.LastNotification = time.Now()

	s.updateNotification(notification)
}

func (s *Scheduler) buildMessage(stream *api.Stream) string {
	if stream.ViewerCount > 1 {
		return fmt.Sprintf(
			"%s is streaming %s with %d viewers.",
			commands.MakeMarkdownLinkUser(stream.UserDisplayName),
			stream.GameName,
			stream.ViewerCount,
		)
	}

	if stream.ViewerCount == 1 {
		return fmt.Sprintf(
			"%s is streaming %s with a single viewer.",
			commands.MakeMarkdownLinkUser(stream.UserDisplayName),
			stream.GameName,
		)
	}

	return fmt.Sprintf(
		"%s just started streaming %s.",
		commands.MakeMarkdownLinkUser(stream.UserDisplayName),
		stream.GameName,
	)
}

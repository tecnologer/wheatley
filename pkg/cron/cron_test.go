package cron //nolint:testpackage // This package is internal and it's being tested by the scheduler_test.go file.

import (
	"context"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/telegram"
	"github.com/tecnologer/wheatley/pkg/twitch"
)

func TestScheduler_NewScheduler(t *testing.T) {
	t.Parallel()

	scheduler, err := NewScheduler(testConfig(t))
	require.NoError(t, err)
	require.NotNil(t, scheduler)
}

func TestScheduler_requireSendMessage(t *testing.T) { //nolint:funlen
	t.Parallel()

	scheduler := &Scheduler{
		Config: &Config{
			NotificationDelay: 5 * time.Minute,
		},
	}

	tests := []struct {
		name         string
		notification *models.Notification
		game         string
		want         bool
	}{
		{
			name: "no_last_notification",
			notification: &models.Notification{
				LastGame: "GameName",
			},
			game: "GameName",
			want: true,
		},
		{
			name: "same_game_outdated_notification",
			notification: &models.Notification{
				LastGame:         "GameName",
				LastNotification: time.Now().Add(-6 * time.Minute),
			},
			game: "GameName",
			want: true,
		},
		{
			name: "same_game_recent_notification",
			notification: &models.Notification{
				LastGame:         "GameName",
				LastNotification: time.Now().Add(-4 * time.Minute),
			},
			game: "GameName",
			want: false,
		},
		{
			name: "different_game_outdated_notification",
			notification: &models.Notification{
				LastGame:         "GameName",
				LastNotification: time.Now().Add(-6 * time.Minute),
			},
			game: "AnotherGame",
			want: true,
		},
		{
			name: "different_game_recent_notification",
			notification: &models.Notification{
				LastGame:         "GameName",
				LastNotification: time.Now().Add(-4 * time.Minute),
			},
			game: "AnotherGame",
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := scheduler.requireSendMessage(test.notification, test.game)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestScheduler_buildMessage(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name         string
		stream       *api.Stream
		notification *models.Notification
		want         string
	}{
		{
			name: "single_viewer",
			stream: &api.Stream{
				ViewerCount:     1,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{},
			want:         "[StreamerName](https://twitch.tv/StreamerName) is now streaming GameName with a single viewer.",
		},
		{
			name: "multiple_viewers",
			stream: &api.Stream{
				ViewerCount:     6,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{},
			want:         "[StreamerName](https://twitch.tv/StreamerName) is now streaming GameName with 6 viewers.",
		},
		{
			name: "single_viewer_different_game",
			stream: &api.Stream{
				ViewerCount:     1,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{
				LastGame:         "AnotherGame",
				LastNotification: time.Now().Add(-6 * time.Minute),
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) changed the game from AnotherGame to GameName with a single viewer.",
		},
		{
			name: "multiple_viewers_different_game",
			stream: &api.Stream{
				ViewerCount:     6,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{
				LastGame:         "AnotherGame",
				LastNotification: time.Now().Add(-6 * time.Minute),
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) changed the game from AnotherGame to GameName with 6 viewers.",
		},
		{
			name: "no_viewers",
			stream: &api.Stream{
				ViewerCount:     0,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{},
			want:         "[StreamerName](https://twitch.tv/StreamerName) just started streaming GameName.",
		},
		{
			name: "no_viewers_different_game",
			stream: &api.Stream{
				ViewerCount:     0,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			notification: &models.Notification{
				LastGame:         "AnotherGame",
				LastNotification: time.Now().Add(-6 * time.Minute),
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) changed the game from AnotherGame to GameName with no viewers.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := &Scheduler{}
			got := s.buildMessage(test.stream, test.notification)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestScheduler_notifyStreamerWentOffline(t *testing.T) {
	t.Parallel()

	scheduler, err := NewScheduler(testConfig(t))
	require.NoError(t, err)

	now := time.Now()

	notification := &models.Notification{
		TwitchStreamerName: "StreamerName",
		TelegramChatID:     123,
		LastNotification:   now,
	}

	assert.Equal(t, now, notification.LastNotification)

	scheduler.notifyStreamerWentOffline(notification)
	assert.Equal(t, time.Time{}, notification.LastNotification)
}

func testConfig(t *testing.T) *Config {
	t.Helper()

	tBot, err := telegram.NewBot(&telegram.Config{
		IsMock: true,
		Token:  "telegram_token",
		TwitchConfig: &twitch.Config{
			ClientID:     "twitch_client",
			ClientSecret: "twitch_secret",
		},
	})
	require.NoError(t, err, "create telegram bot")

	return &Config{
		SchedulerInterval: 5 * time.Minute,
		NotificationDelay: 2 * time.Hour,
		Context:           context.Background(),
		TwitchConfig: &twitch.Config{
			ClientID:     "twitch_client",
			ClientSecret: "twitch_secret",
		},
		TelegramBot:   tBot,
		Notifications: dao.NewNotificationsDAO(nil),
	}
}

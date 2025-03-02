package telegram_test

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/dao/db"
	"github.com/tecnologer/wheatley/pkg/telegram"
)

func TestBot_NewBot(t *testing.T) {
	t.Parallel()

	t.Run("incomplete_config", func(t *testing.T) {
		t.Parallel()

		bot, err := telegram.NewBot(&telegram.Config{
			IsMock: true,
		})
		require.Error(t, err)
		require.Nil(t, bot)
	})

	t.Run("valid_config", func(t *testing.T) {
		t.Parallel()

		bot, err := telegram.NewBot(&telegram.Config{
			IsMock: true,
			Token:  "test_token",
			DB:     &db.Connection{},
		})

		require.NoError(t, err)
		require.NotNil(t, bot)
	})
}

func TestBot_SendMessage(t *testing.T) {
	t.Parallel()

	bot, err := telegram.NewBot(&telegram.Config{
		IsMock: true,
		Token:  "test_token",
		DB:     &db.Connection{},
	})
	require.NoError(t, err)
	require.NotNil(t, bot)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := bot.SendMessage(0, "test")
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestBot_ReadUpdates(t *testing.T) {
	t.Parallel()

	bot, err := telegram.NewBot(&telegram.Config{
		IsMock: true,
		Token:  "test_token",
		DB:     nil,
	})
	require.NoError(t, err)
	require.NotNil(t, bot)

	done := make(chan bool)

	go func() {
		_ = bot.ReadUpdates()

		done <- true
	}()

	time.Sleep(time.Millisecond * 500)

	mockBot, ok := bot.Bot.(*telegram.Mock)
	require.True(t, ok)

	mockBot.SendUpdates(messageUpdateCmdAdd(t))
	mockBot.CloseUpdates()

	<-done
}

func TestBot_NotifyAdminIfNecessary(t *testing.T) {
	t.Parallel()

	bot, err := telegram.NewBot(&telegram.Config{
		IsMock:    true,
		Token:     "test_token",
		DB:        nil,
		Admins:    []string{"admin"},
		ChatAdmin: 123,
	})
	require.NoError(t, err)
	require.NotNil(t, bot)

	t.Run("no_need_to_notify", func(t *testing.T) {
		t.Parallel()

		bot.NotifyAdminIfNecessary(messageUpdateCmdStart(t))
	})

	t.Run("need_to_notify", func(t *testing.T) {
		t.Parallel()

		bot.NotifyAdminIfNecessary(messageUpdateCmdAdd(t))
	})

}

func messageUpdateCmdAdd(t *testing.T) tgbotapi.Update {
	t.Helper()

	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add streamer_name",
			From: &tgbotapi.User{
				ID:       123,
				UserName: "user_name",
			},
			Chat: &tgbotapi.Chat{
				ID:       123,
				UserName: "chat_name",
			},
		},
	}
}

func messageUpdateCmdStart(t *testing.T) tgbotapi.Update {
	t.Helper()

	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			From: &tgbotapi.User{
				ID:       123,
				UserName: "user_name",
			},
			Chat: &tgbotapi.Chat{
				ID:       123,
				UserName: "chat_name",
			},
		},
	}
}

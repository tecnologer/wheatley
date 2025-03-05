package cli //nolint:testpackage

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/cmd/flags"
	"github.com/tecnologer/wheatley/pkg/cron"
	"github.com/tecnologer/wheatley/pkg/dao"
	"github.com/tecnologer/wheatley/pkg/twitch"
	"github.com/urfave/cli/v2"
)

func TestNewCLI(t *testing.T) {
	t.Parallel()

	newCLI := NewCLI("dev")
	require.NotNil(t, newCLI)

	t.Run("flags", func(t *testing.T) {
		t.Parallel()

		want := []string{
			flags.VerboseFlagName,
			flags.VerboseFlagAlias,
			flags.DBPasswordFlagName,
			flags.DBPasswordFlagAlias,
			flags.DBNameFlagName,
			flags.DBNameFlagAlias,
			flags.TelegramTokenFlagName,
			flags.TelegramTokenFlagAlias,
			flags.IntervalFlagName,
			flags.IntervalFlagAlias,
			flags.ResendIntervalFlagName,
			flags.ResendIntervalFlagAlias,
			flags.TwitchClientIDFlagName,
			flags.TwitchClientIDFlagAlias,
			flags.TwitchClientSecretFlagName,
			flags.TwitchClientSecretFlagAlias,
			flags.TelegramAdminChatIDsFlagName,
			flags.TelegramAdminChatIDsFlagAlias,
			flags.TelegramAdminsFlagName,
			flags.TelegramAdminsFlagAlias,
		}

		for _, cliFlag := range newCLI.Flags {
			for _, flagName := range cliFlag.Names() {
				assert.Contains(t, want, flagName)
			}
		}
	})
}

func TestApp_beforeRun(t *testing.T) {
	t.Parallel()

	newCLI := NewCLI("dev")
	require.NotNil(t, newCLI)

	tests := []struct {
		name          string
		verbose       bool
		telegramToken string
		dbName        string
		wantErr       error
	}{
		{
			name:          "no_telegram_token",
			verbose:       false,
			telegramToken: "",
			dbName:        "test",
			wantErr:       fmt.Errorf("telegram token is required"),
		},
		{
			name:          "no_db_name",
			verbose:       false,
			telegramToken: "test",
			dbName:        "",
			wantErr:       fmt.Errorf("db name is required"),
		},
		{
			name:          "success_no_verbose",
			verbose:       false,
			telegramToken: "test",
			dbName:        "test",
			wantErr:       nil,
		},
		{
			name:          "success_verbose",
			verbose:       true,
			telegramToken: "test",
			dbName:        "test",
			wantErr:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flagSet.String(flags.TelegramTokenFlagName, test.telegramToken, "")
			flagSet.String(flags.DBNameFlagName, test.dbName, "")
			flagSet.Bool(flags.VerboseFlagName, test.verbose, "")

			ctx := cli.NewContext(newCLI.App, flagSet, nil)

			err := newCLI.beforeRun(ctx)
			require.Equal(t, test.wantErr, err)
		})
	}
}

func TestApp_createConnection(t *testing.T) {
	t.Cleanup(func() {
		err := os.Remove("./test.db")
		require.NoError(t, err, "remove test.db")
	})

	t.Parallel()

	newCLI := NewCLI("dev")
	require.NotNil(t, newCLI)

	tests := []struct {
		name    string
		dbName  string
		wantErr error
	}{
		{
			name:    "no_db_name",
			dbName:  "",
			wantErr: fmt.Errorf("create new connection: create DB file name: invalid db config: missing db name"),
		},
		{
			name:    "success",
			dbName:  "test",
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flagSet.String(flags.DBNameFlagName, test.dbName, "")

			ctx := cli.NewContext(newCLI.App, flagSet, nil)

			_, err := newCLI.createConnection(ctx)
			if test.wantErr == nil {
				require.NoError(t, err)
				return
			}

			assert.Equal(t, test.wantErr.Error(), err.Error())
		})
	}
}

func TestApp_cronOptions(t *testing.T) {
	t.Parallel()

	newCLI := NewCLI("dev")
	require.NotNil(t, newCLI)

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Bool(flags.VerboseFlagName, true, "")
	flagSet.String(flags.DBPasswordFlagName, "test_pwd", "")
	flagSet.String(flags.DBNameFlagName, "db_user", "")
	flagSet.String(flags.TelegramTokenFlagName, "telegram_token", "")
	flagSet.Int(flags.IntervalFlagName, 5, "")
	flagSet.Int(flags.ResendIntervalFlagName, 2, "")
	flagSet.String(flags.TwitchClientIDFlagName, "twitch_client", "")
	flagSet.String(flags.TwitchClientSecretFlagName, "twitch_secret", "")
	flagSet.Int(flags.TelegramAdminChatIDsFlagName, 1234, "")
	flagSet.String(flags.TelegramAdminsFlagName, "tester", "")

	want := &cron.Config{
		SchedulerInterval: 5 * time.Minute,
		NotificationDelay: 2 * time.Hour,
		Context:           context.Background(),
		TwitchConfig: &twitch.Config{
			ClientID:     "twitch_client",
			ClientSecret: "twitch_secret",
		},
		TelegramBot:   newCLI.bot,
		Notifications: dao.NewNotifications(nil),
	}

	ctx := cli.NewContext(newCLI.App, flagSet, nil)

	config := newCLI.cronOptions(context.Background(), ctx, nil)
	require.NotNil(t, config)
	assert.Equal(t, want, config)
}

package commands_test

import (
	"context"
	"testing"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
	"github.com/tecnologer/wheatley/pkg/twitch"
)

func TestAddStreamerCmd(t *testing.T) { //nolint:funlen
	t.Parallel()

	cmd := commands.AddStreamerCmd(nil)
	require.NotNil(t, cmd)
	require.NotNil(t, cmd.Handler)
	require.NotNil(t, cmd.Help)

	assert.Equal(t, commands.AddStreamerCmdName, cmd.Name)
	assert.Equal(t, "Adds a streamer to the list of notifications. You will be notified when the streamer goes live.", cmd.Description)
	assert.Equal(t, "Usage: `/add <streamer_name>`", cmd.Help())

	t.Run("add_streamer_handler_message", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t), "streamer_name")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(t, "Done! I'll notify you when `streamer_name` goes live", response.Message(), "response message should match")
	})

	t.Run("add_streamer_handler_message_missing_streamer_name", func(t *testing.T) {
		t.Parallel()

		update := messageUpdateCmdAdd(t)
		update.Message.Text = "/add"

		response := cmd.Handler(cmd, update)
		require.NotNil(t, response, "response should not be nil")
		require.True(t, response.HasError(), "response should have an error")

		assert.Equal(
			t,
			"❌ Error adding streamer: missing streamer name\n\nUsage: `/add <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("add_streamer_handler_edit_message", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, editMessageUpdateCmdAdd(t), "streamer_name")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(t, "Done! I'll notify you when `streamer_name` goes live", response.Message(), "response message should match")
	})

	t.Run("add_streamer_handler_message_missing_streamer_name", func(t *testing.T) {
		t.Parallel()

		update := editMessageUpdateCmdAdd(t)
		update.EditedMessage.Text = "/add"

		response := cmd.Handler(cmd, update)
		require.NotNil(t, response, "response should not be nil")
		require.True(t, response.HasError(), "response should have an error")

		assert.Equal(
			t,
			"❌ Error adding streamer: missing streamer name\n\nUsage: `/add <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})
}

func TestRemoveStreamerCmd(t *testing.T) { //nolint:funlen
	t.Parallel()

	cmd := commands.RemoveStreamerCmd(nil)
	require.NotNil(t, cmd)
	require.NotNil(t, cmd.Handler)
	require.NotNil(t, cmd.Help)

	assert.Equal(t, commands.RemoveStreamerCmdName, cmd.Name)
	assert.Equal(t, "Removes a streamer from the list. You won't receive notifications for this streamer anymore.", cmd.Description)
	assert.Equal(t, "Usage: `/remove <streamer_name>`", cmd.Help())

	t.Run("remove_streamer_handler_message", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t), "streamer_name")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(t, "Done! You won't receive notifications for `streamer_name` anymore", response.Message(), "response message should match")
	})

	t.Run("remove_streamer_handler_message_missing_streamer_name", func(t *testing.T) {
		t.Parallel()

		update := messageUpdateCmdAdd(t)
		update.Message.Text = "/remove"

		response := cmd.Handler(cmd, update)
		require.NotNil(t, response, "response should not be nil")
		require.True(t, response.HasError(), "response should have an error")

		assert.Equal(
			t,
			"❌ Error removing streamer: missing streamer name\n\nUsage: `/remove <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("remove_streamer_handler_edit_message", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, editMessageUpdateCmdAdd(t), "streamer_name")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(t, "Done! You won't receive notifications for `streamer_name` anymore", response.Message(), "response message should match")
	})

	t.Run("remove_streamer_handler_edit_message_missing_streamer_name", func(t *testing.T) {
		t.Parallel()

		update := editMessageUpdateCmdAdd(t)
		update.EditedMessage.Text = "/remove"

		response := cmd.Handler(cmd, update)
		require.NotNil(t, response, "response should not be nil")
		require.True(t, response.HasError(), "response should have an error")

		assert.Equal(
			t,
			"❌ Error removing streamer: missing streamer name\n\nUsage: `/remove <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})
}

func TestListStreamersCmd(t *testing.T) {
	t.Parallel()

	twitchInstance := twitch.New(&twitch.Config{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		IsMock:       true,
	})

	twitchMock, isMock := twitchInstance.(*twitch.MockAPI)
	require.True(t, isMock)

	cmd := commands.ListStreamersCmd(nil, twitchInstance)
	require.NotNil(t, cmd)
	require.Nil(t, cmd.Help)

	assert.Equal(t, commands.ListStreamersCmdName, cmd.Name)
	assert.Equal(t, "Lists all the streamers you're currently following.", cmd.Description)

	t.Run("list_streamers_handler", func(t *testing.T) {
		t.Parallel()

		twitchMock.On("StreamByName", context.Background(), "streamer_name").Return(&api.Stream{
			GameName: "Game Name",
		}, nil)

		twitchMock.On("StreamByName", context.Background(), "another_streamer").Return(nil, nil)

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t))
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"You are subscribed to the following streamers:\n\n"+
				"🎮 [streamer_name](https://twitch.tv/streamer_name) - (Playing: Game Name)\n"+
				"🎮 [another_streamer](https://twitch.tv/another_streamer)  - (offline)\n",
			response.Message(),
			"response message should match",
		)
	})
}

func TestStartCmd(t *testing.T) {
	t.Parallel()

	cmd := commands.StartCmd()
	require.NotNil(t, cmd)
	require.Nil(t, cmd.Help)

	assert.Equal(t, commands.StartCmdName, cmd.Name)
	assert.Equal(t, "Starts the bot.", cmd.Description)

	t.Run("start_handler", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t))
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"Hello! I'm Wheatley, your Twitch notifications bot. To add a streamer to the list of notifications, "+
				"use the command `/add <streamer_name>` or /help.\n\n Source code on [GitHub](https://github.com/tecnologer/wheatley), "+
				"Need more help? Contact @tecnologer",
			response.Message(),
			"response message should match",
		)
	})
}

func TestHelpCmd(t *testing.T) { //nolint:funlen
	t.Parallel()

	allCommands := commands.NewCommands(nil, twitch.New(&twitch.Config{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		IsMock:       true,
	}))

	cmd := commands.HelpCmd(allCommands)
	require.NotNil(t, cmd)
	require.NotNil(t, cmd.Handler)
	require.NotNil(t, cmd.Help)

	assert.Equal(t, commands.HelpCmdName, cmd.Name)
	assert.Equal(t, "Shows the available commands and their usage.", cmd.Description)
	assert.Equal(t, "Usage: \n- All available commands /help.\n- Help for specific command `/help <command>`. For example: `/help add`", cmd.Help())

	t.Run("help_handler", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t))
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"Hello! I'm Wheatley, your Twitch notifications bot, and I can help you with the following commands:\n\n"+
				"» /add: Adds a streamer to the list of notifications. You will be notified when the streamer goes live.\n"+
				"» /help: Shows the available commands and their usage.\n"+
				"» /list: Lists all the streamers you're currently following.\n"+
				"» /remove: Removes a streamer from the list. You won't receive notifications for this streamer anymore.\n"+
				"» /start: Starts the bot.\n\n"+
				"If you want get more information about a specific command, use `/help <command>`. For example: `/help add`.",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("help_handler_edit", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, editMessageUpdateCmdAdd(t))
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"Hello! I'm Wheatley, your Twitch notifications bot, and I can help you with the following commands:\n\n"+
				"» /add: Adds a streamer to the list of notifications. You will be notified when the streamer goes live.\n"+
				"» /help: Shows the available commands and their usage.\n"+
				"» /list: Lists all the streamers you're currently following.\n"+
				"» /remove: Removes a streamer from the list. You won't receive notifications for this streamer anymore.\n"+
				"» /start: Starts the bot.\n\n"+
				"If you want get more information about a specific command, use `/help <command>`. For example: `/help add`.",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("help_for_command", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t), "add")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"/add: Adds a streamer to the list of notifications. You will be notified when the streamer goes live.\n\n"+
				"Usage: `/add <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("help_for_command_edited_message", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, editMessageUpdateCmdAdd(t), "add")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"/add: Adds a streamer to the list of notifications. You will be notified when the streamer goes live.\n\n"+
				"Usage: `/add <streamer_name>`",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("help_for_command_not_found", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t), "not_found")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should have an error")

		assert.Equal(
			t,
			"Command not_found not found. Available commands:\n\n"+
				"Hello! I'm Wheatley, your Twitch notifications bot, and I can help you with the following commands:\n\n"+
				"» /add: Adds a streamer to the list of notifications. You will be notified when the streamer goes live.\n"+
				"» /help: Shows the available commands and their usage.\n» /list: Lists all the streamers you're currently following.\n"+
				"» /remove: Removes a streamer from the list. You won't receive notifications for this streamer anymore.\n"+
				"» /start: Starts the bot.\n\nIf you want get more information about a specific command, use `/help <command>`. For example: `/help add`.",
			response.Message(),
			"response message should match",
		)
	})

	t.Run("help_command_not_help_handler", func(t *testing.T) {
		t.Parallel()

		response := cmd.Handler(cmd, messageUpdateCmdAdd(t), "start")
		require.NotNil(t, response, "response should not be nil")
		require.False(t, response.HasError(), "response should not have an error")

		assert.Equal(
			t,
			"/start: Starts the bot.",
			response.Message(),
			"response message should match",
		)
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
			Chat: tgbotapi.Chat{
				ID:       123,
				UserName: "chat_name",
			},
		},
	}
}

func editMessageUpdateCmdAdd(t *testing.T) tgbotapi.Update {
	t.Helper()

	return tgbotapi.Update{
		EditedMessage: &tgbotapi.Message{
			Text: "/add streamer_name",
			Chat: tgbotapi.Chat{
				ID:       123,
				UserName: "chat_name",
			},
		},
	}
}

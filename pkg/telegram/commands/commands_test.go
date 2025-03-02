package commands_test

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
)

func TestCommands_HasHandler(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName commands.CommandName
		want    bool
	}{
		{
			name:    "start_cmd",
			cmdName: commands.StartCmdName,
			want:    true,
		},
		{
			name:    "add_streamer_cmd",
			cmdName: commands.AddStreamerCmdName,
			want:    true,
		},
		{
			name:    "remove_streamer_cmd",
			cmdName: commands.RemoveStreamerCmdName,
			want:    true,
		},
		{
			name:    "help_cmd",
			cmdName: commands.HelpCmdName,
			want:    true,
		},
		{
			name:    "list_streamers_cmd",
			cmdName: commands.ListStreamersCmdName,
			want:    true,
		},
		{
			name:    "non_existent_cmd",
			cmdName: "non_existent_cmd",
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.HasHandler(test.cmdName)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCommands_HasHelp(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName commands.CommandName
		want    bool
	}{
		{
			name:    "start_cmd",
			cmdName: commands.StartCmdName,
			want:    false,
		},
		{
			name:    "add_streamer_cmd",
			cmdName: commands.AddStreamerCmdName,
			want:    true,
		},
		{
			name:    "remove_streamer_cmd",
			cmdName: commands.RemoveStreamerCmdName,
			want:    true,
		},
		{
			name:    "help_cmd",
			cmdName: commands.HelpCmdName,
			want:    true,
		},
		{
			name:    "list_streamers_cmd",
			cmdName: commands.ListStreamersCmdName,
			want:    false,
		},
		{
			name:    "non_existent_cmd",
			cmdName: "non_existent_cmd",
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.HasHelp(test.cmdName)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCommands_Execute(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName string
		update  tgbotapi.Update
		args    []string
		want    string
	}{
		{
			name:    "start_cmd",
			cmdName: string(commands.StartCmdName),
			update:  messageUpdateCmdAdd(t),
			args:    nil,
			want: "Hello! I'm Wheatley, your Twitch notifications bot. To add a streamer to the list of notifications, use the command " +
				"`/add <streamer_name>` or /help.\n\n Source code on [GitHub](https://github.com/tecnologer/wheatley), " +
				"Need more help? Contact @tecnologer",
		},
		{
			name:    "unknown_cmd",
			cmdName: "unknown_cmd",
			update:  messageUpdateCmdAdd(t),
			args:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.Execute(test.cmdName, test.update, test.args...)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCommand_Help(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName string
		want    string
	}{
		{
			name:    "no_help_handler",
			cmdName: string(commands.StartCmdName),
			want:    "",
		},
		{
			name:    "help_handler",
			cmdName: string(commands.AddStreamerCmdName),
			want:    "Usage: `/add <streamer_name>`",
		},
		{
			name:    "unknown_cmd",
			cmdName: "unknown_cmd",
			want:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.Help(test.cmdName)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCommand_HasAdminNotification(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName string
		want    bool
	}{
		{
			name:    "no_notification_handler",
			cmdName: string(commands.StartCmdName),
			want:    false,
		},
		{
			name:    "notification_handler",
			cmdName: string(commands.AddStreamerCmdName),
			want:    true,
		},
		{
			name:    "unknown_cmd",
			cmdName: "unknown_cmd",
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.HasAdminNotification(test.cmdName)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCommand_AdminNotification(t *testing.T) {
	t.Parallel()

	allCommands := commands.NewCommands(nil)

	tests := []struct {
		name    string
		cmdName string
		update  tgbotapi.Update
		args    []string
		want    *commands.Response
	}{
		{
			name:    "no_notification_handler",
			cmdName: string(commands.StartCmdName),
			want:    nil,
		},
		{
			name:    "notification_handler_add",
			cmdName: string(commands.AddStreamerCmdName),
			update:  messageUpdateCmdAdd(t),
			args:    []string{"streamer_name"},
			want: commands.NewResponse(
				commands.WithCommand(allCommands.Map[commands.AddStreamerCmdName]),
				commands.WithMessage("The user @user_name added the streamer `streamer_name` to chat \"chat_name\"."),
			),
		},
		{
			name:    "notification_handler_remove",
			cmdName: string(commands.RemoveStreamerCmdName),
			update:  messageUpdateCmdAdd(t),
			args:    []string{"streamer_name"},
			want: commands.NewResponse(
				commands.WithCommand(allCommands.Map[commands.RemoveStreamerCmdName]),
				commands.WithMessage("The user @user_name removed the streamer `streamer_name` from chat \"chat_name\"."),
			),
		},
		{
			name:    "unknown_cmd",
			cmdName: "unknown_cmd",
			want:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := allCommands.AdminNotification(test.cmdName, test.update, test.args...)
			assert.Equal(t, test.want, got)
		})
	}
}

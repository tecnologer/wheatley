package message_test

import (
	"testing"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

func TestGetFromUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   string
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			want:   "",
		},
		{
			name: "empty_message_no_nil_message_nil",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{},
			},
			want: "",
		},
		{
			name: "from_message",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: "message_text",
				},
			},
			want: "message_text",
		},
		{
			name: "from_edited_message",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					Text: "edited_message_text",
				},
			},
			want: "edited_message_text",
		},
		{
			name: "message_and_edited_message_not_nil",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: "message_text",
				},
				EditedMessage: &tgbotapi.Message{
					Text: "edited_message_text",
				},
			},
			want: "message_text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := message.GetFromUpdate(tt.update)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetChatIDFromUpdate(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   int64
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			want:   0,
		},
		{
			name: "empty_message_no_nil_message_nil",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{},
			},
			want: 0,
		},
		{
			name: "not_chat_message_or_edited_message",
			update: tgbotapi.Update{
				Message:       &tgbotapi.Message{},
				EditedMessage: &tgbotapi.Message{},
			},
			want: 0,
		},
		{
			name: "from_message",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						ID: 123,
					},
				},
			},
			want: 123,
		},
		{
			name: "from_edited_message",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						ID: 123,
					},
				},
			},
			want: 123,
		},
		{
			name: "message_and_edited_message_not_nil",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						ID: 123,
					},
				},
				EditedMessage: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						ID: 456,
					},
				},
			},
			want: 123,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := message.GetChatIDFromUpdate(test.update)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetChatNameFromUpdate(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   string
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			want:   "<<no defined>>",
		},
		{
			name: "empty_message_edited_message_not_nil",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{},
				Message:       &tgbotapi.Message{},
			},
			want: "",
		},
		{
			name: "from_message_username",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						UserName: "chat_name",
					},
				},
			},
			want: "chat_name",
		},
		{
			name: "from_message_title",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						Title: "chat_title",
					},
				},
			},
			want: "chat_title",
		},
		{
			name: "from_edited_message_username",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						UserName: "chat_name",
					},
				},
			},
			want: "chat_name",
		},
		{
			name: "from_edited_message_title",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					Chat: tgbotapi.Chat{
						Title: "chat_title",
					},
				},
			},
			want: "chat_title",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := message.GetChatNameFromUpdate(test.update)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIsBot(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   bool
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			want:   false,
		},
		{
			name: "empty_message_edited_message_not_nil",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{},
				Message:       &tgbotapi.Message{},
			},
			want: false,
		},
		{
			name: "from_message_is_bot",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: true,
					},
				},
			},
			want: true,
		},
		{
			name: "from_message_is_not_bot",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: false,
					},
				},
			},
			want: false,
		},
		{
			name: "from_edited_message_is_bot",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: true,
					},
				},
			},
			want: true,
		},
		{
			name: "from_edited_message_is_not_bot",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: false,
					},
				},
			},
			want: false,
		},
		{
			name: "from_message_and_edited_message_is_bot",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: true,
					},
				},
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						IsBot: false,
					},
				},
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := message.IsBot(test.update)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestSenderName(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   string
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			want:   "",
		},
		{
			name: "empty_message_edited_message_not_nil",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{},
				Message:       &tgbotapi.Message{},
			},
			want: "",
		},
		{
			name: "from_message_username",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "user_name",
					},
				},
			},
			want: "user_name",
		},
		{
			name: "from_edited_message_username",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "user_name",
					},
				},
			},
			want: "user_name",
		},
		{
			name: "from_message_and_edited_message_username",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "user_name",
					},
				},
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "edited_user_name",
					},
				},
			},
			want: "user_name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := message.SenderName(test.update)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestSentByAdmin(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		admins []string
		want   bool
	}{
		{
			name:   "empty",
			update: tgbotapi.Update{},
			admins: []string{"admin"},
			want:   false,
		},
		{
			name: "message_is_admin",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "admin",
					},
				},
			},
			admins: []string{"admin"},
			want:   true,
		},
		{
			name: "edited_message_is_admin",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "admin",
					},
				},
			},
			admins: []string{"admin"},
			want:   true,
		},
		{
			name: "message_and_edited_message_is_admin",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "admin",
					},
				},
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "admin",
					},
				},
			},
			admins: []string{"admin"},
			want:   true,
		},
		{
			name: "message_is_not_admin",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "user",
					},
				},
			},
			admins: []string{"admin"},
			want:   false,
		},
		{
			name: "edited_message_is_not_admin",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					From: &tgbotapi.User{
						UserName: "user",
					},
				},
			},
			admins: []string{"admin"},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := message.SentByAdmin(test.update, test.admins)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetMessageThreadID(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name   string
		update tgbotapi.Update
		want   int
	}{
		{
			name:   "empty_update",
			update: tgbotapi.Update{},
			want:   0,
		},
		{
			name: "empty_message",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{},
			},
			want: 0,
		},
		{
			name: "message_with_empty_reply",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					ReplyToMessage: nil,
				},
			},
			want: 0,
		},
		{
			name: "message_with_reply",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					ReplyToMessage: &tgbotapi.Message{
						MessageThreadID: 123,
					},
				},
			},
			want: 123,
		},
		{
			name: "edited_message_with_reply",
			update: tgbotapi.Update{
				EditedMessage: &tgbotapi.Message{
					ReplyToMessage: &tgbotapi.Message{
						MessageThreadID: 456,
					},
				},
			},
			want: 456,
		},
		{
			name: "message_and_edited_message_with_reply",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					ReplyToMessage: &tgbotapi.Message{
						MessageThreadID: 123,
					},
				},
				EditedMessage: &tgbotapi.Message{
					ReplyToMessage: &tgbotapi.Message{
						MessageThreadID: 456,
					},
				},
			},
			want: 123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := message.GetMessageThreadID(tt.update)
			assert.Equal(t, tt.want, got)
		})
	}
}

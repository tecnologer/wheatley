package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/wheatley/pkg/telegram/commands"
)

func TestNewResponse(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name        string
		response    *commands.Response
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "empty_response",
			response:    commands.NewResponse(),
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:        "response_with_message",
			response:    commands.NewResponse().SetMessage("test message"),
			wantMessage: "test message",
			wantErr:     false,
		},
		{
			name:        "response_with_error",
			response:    commands.NewResponse().SetError("test error"),
			wantMessage: "❌ test error",
			wantErr:     true,
		},
		{
			name:        "response_with_message_and_error",
			response:    commands.NewResponse().SetMessage("test message").SetError("test error"),
			wantMessage: "❌ test error",
			wantErr:     true,
		},
		{
			name: "response_with_all_data",
			response: commands.NewResponse(
				commands.WithMessage("test message"),
				commands.WithMissingArgs(),
				commands.WithIsError(),
			),
			wantMessage: "❌❓ test message",
			wantErr:     true,
		},
		{
			name:        "response_with_messagef",
			response:    commands.NewResponse(commands.WithMessagef("test %s with format", "message")),
			wantMessage: "test message with format",
			wantErr:     false,
		},
		{
			name:        "response_with_missing_args",
			response:    commands.NewResponse().SetMissingArgs("test missing args"),
			wantMessage: "❌❓ test missing args",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.wantMessage, test.response.Message())
			assert.Equal(t, test.wantErr, test.response.HasError())
		})
	}
}

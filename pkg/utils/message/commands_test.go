package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/utils/message"
)

func TestExtractCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		message  string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "empty_message",
			message:  "",
			wantCmd:  "",
			wantArgs: nil,
		},
		{
			name:     "no_command",
			message:  "hello",
			wantCmd:  "",
			wantArgs: nil,
		},
		{
			name:     "command_no_args",
			message:  "/start",
			wantCmd:  "start",
			wantArgs: nil,
		},
		{
			name:     "command_with_args",
			message:  "/echo hello world",
			wantCmd:  "echo",
			wantArgs: []string{"hello", "world"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotCmd, gotArgs := message.ExtractCommand(test.message)
			assert.Equal(t, test.wantCmd, gotCmd)
			assert.Equal(t, test.wantArgs, gotArgs)
		})
	}
}

func TestExtractValueFromArg(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name         string
		arg          string
		wantArgName  string
		wantArgValue string
	}{
		{
			name:         "empty_arg",
			arg:          "",
			wantArgName:  "",
			wantArgValue: "",
		},
		{
			name:         "arg_with_spaces",
			arg:          "key value",
			wantArgName:  "key value",
			wantArgValue: "",
		},
		{
			name:         "arg_with_colon",
			arg:          "key:value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
		{
			name:         "arg_with_equal",
			arg:          "key=value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
		{
			name:         "arg_with_spaces_colon",
			arg:          "key: value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
		{
			name:         "arg_with_spaces_equal",
			arg:          "key = value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
		{
			name:         "arg_with_spaces_colon_spaces",
			arg:          "key : value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
		{
			name:         "arg_with_spaces_equal_spaces",
			arg:          "key = value",
			wantArgName:  "key",
			wantArgValue: "value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			argName, argValue := message.ExtractValueFromArg(test.arg)
			assert.Equal(t, test.wantArgName, argName)
			assert.Equal(t, test.wantArgValue, argValue)
		})
	}
}

func TestArgsToMap(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		order   []string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "empty_args",
			args:    nil,
			order:   nil,
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "empty_order",
			args:    []string{"key value"},
			order:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:  "args_no_named",
			args:  []string{"value", "value2"},
			order: []string{"key", "key2"},
			want: map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		},
		{
			name:  "args_named",
			args:  []string{"key: value", "key2: value2"},
			order: []string{"key", "key2"},
			want: map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		},
		{
			name:  "args_mixed_named",
			args:  []string{"key: value", "value2"},
			order: []string{"key", "key2"},
			want: map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := message.ArgsToMap(test.args, test.order)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

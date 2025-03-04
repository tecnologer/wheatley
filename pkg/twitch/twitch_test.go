package twitch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tecnologer/wheatley/pkg/twitch"
)

func TestTwitch_StreamByName(t *testing.T) { //nolint: funlen
	t.Parallel()

	tests := []struct {
		name          string
		mockSetup     func(m *twitch.MockAPI)
		username      string
		expectedError string
		want          *api.Stream
	}{
		{
			name:     "successful_stream_by_name",
			username: "streamer_name",
			mockSetup: func(m *twitch.MockAPI) {
				m.On("StreamByName", mock.Anything, mock.Anything).
					Return(&api.Stream{
						UserDisplayName: "Streamer Name",
					}, nil)
			},
			want: &api.Stream{
				UserDisplayName: "Streamer Name",
			},
		},
		{
			name: "streamer_not_found",
			mockSetup: func(m *twitch.MockAPI) {
				m.On("StreamByName", mock.Anything, mock.Anything).
					Return(nil, errors.New("user not found"))
			},
			expectedError: "getting stream: user not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := twitch.New(&twitch.Config{
				ClientID:     "client_id",
				ClientSecret: "client_secret",
				IsMock:       true,
			})
			require.NotNil(t, service)

			mockAPI, ok := service.(*twitch.MockAPI)
			require.True(t, ok)

			test.mockSetup(mockAPI)

			stream, err := service.StreamByName(context.Background(), test.username)

			if test.expectedError != "" {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, test.want, stream)

			mockAPI.AssertExpectations(t)
		})
	}
}

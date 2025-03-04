package cron //nolint:testpackage // This package is internal and it's being tested by the scheduler_test.go file.

import (
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_buildMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		stream *api.Stream
		want   string
	}{
		{
			name: "single_viewer",
			stream: &api.Stream{
				ViewerCount:     1,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) is streaming GameName with a single viewer.",
		},
		{
			name: "multiple_viewers",
			stream: &api.Stream{
				ViewerCount:     6,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) is streaming GameName with 6 viewers.",
		},
		{
			name: "no_viewers",
			stream: &api.Stream{
				ViewerCount:     0,
				UserDisplayName: "StreamerName",
				GameName:        "GameName",
			},
			want: "[StreamerName](https://twitch.tv/StreamerName) just started streaming GameName.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := &Scheduler{}
			got := s.buildMessage(test.stream)
			assert.Equal(t, test.want, got)
		})
	}
}

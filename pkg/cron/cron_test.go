package cron

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/tecnologer/wheatley/pkg/twitch"
	"testing"
)

func TestScheduler_taskTwitchCheckStreamers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Scheduler: tt.fields.Scheduler,
				Config:    tt.fields.Config,
				twitch:    tt.fields.twitch,
			}
			s.taskTwitchCheckStreamers()
		})
	}
}

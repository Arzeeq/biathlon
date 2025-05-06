package event_test

import (
	"biathlon/internal/event"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		desc          string
		expectedEvent *event.Event
		isErrExpected bool
	}{
		{
			name: "valid event without extra params",
			desc: "[09:05:59.867] 1 1",
			expectedEvent: &event.Event{
				Time:       time.Date(0, 1, 1, 9, 5, 59, 867000000, time.UTC),
				ID:         1,
				Competitor: 1,
				Extra:      "",
			},
			isErrExpected: false,
		},
		{
			name: "valid event with extra params",
			desc: "[09:15:00.841] 2 1 09:30:00.000",
			expectedEvent: &event.Event{
				Time:       time.Date(0, 1, 1, 9, 15, 0, 841000000, time.UTC),
				ID:         2,
				Competitor: 1,
				Extra:      "09:30:00.000",
			},
			isErrExpected: false,
		},
		{
			name: "valid event with multiple extra params",
			desc: "[09:49:33.123] 6 1 1 2 3",
			expectedEvent: &event.Event{
				Time:       time.Date(0, 1, 1, 9, 49, 33, 123000000, time.UTC),
				ID:         6,
				Competitor: 1,
				Extra:      "1 2 3",
			},
			isErrExpected: false,
		},
		{
			name:          "empty string",
			desc:          "",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "not enough params",
			desc:          "[09:05:59.867] 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid time format",
			desc:          "[09:05:59] 1 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid event ID",
			desc:          "[09:05:59.867] abc 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid competitor ID",
			desc:          "[09:05:59.867] 1 abc",
			expectedEvent: nil,
			isErrExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := event.New(tt.desc)

			if tt.isErrExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedEvent, got)
		})
	}
}

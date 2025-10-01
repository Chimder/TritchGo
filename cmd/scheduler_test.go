package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextInterval(t *testing.T) {

	test := []struct {
		name     string
		interval int
		now      time.Time
		expected time.Duration
	}{
		{
			name:     "15 interval at 2min",
			interval: 15,
			now:      time.Date(2024, 1, 1, 12, 2, 0, 0, time.UTC),
			expected: 13 * time.Minute,
		},
		{
			name:     "15 interval at 25min",
			interval: 15,
			now:      time.Date(2024, 1, 1, 12, 25, 0, 0, time.UTC),
			expected: 5 * time.Minute,
		},
		{
			name:     "15 on 15",
			interval: 15,
			now:      time.Date(2024, 1, 1, 12, 15, 0, 0, time.UTC),
			expected: 15 * time.Minute,
		},
	}

	for _, tt := range test {

		t.Run(tt.name, func(t *testing.T) {

			res := NextIntervalAt(tt.now, tt.interval)
			assert.NotEmpty(t, res)
			assert.Equal(t, tt.expected, res)

		})
	}

}

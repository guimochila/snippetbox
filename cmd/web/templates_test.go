package main

import (
	"testing"
	"time"

	"snippetbox.guimochila.com/internal/assert"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 04, 29, 8, 43, 0, 0, time.UTC),
			want: "29 Apr 2024 at 08:43",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 04, 29, 8, 43, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "29 Apr 2024 at 07:43",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}

}

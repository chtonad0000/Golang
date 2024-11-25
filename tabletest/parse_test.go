package tabletest

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		err      bool
	}{
		{"1s", time.Second, false},
		{"2h45m", 2*time.Hour + 45*time.Minute, false},
		{"-1.5h", -1*time.Hour - 30*time.Minute, false},
		{"0", 0, false},
		{"1.5h", time.Hour + 30*time.Minute, false},
		{"2.5s", 2*time.Second + 500*time.Millisecond, false},
		{"1m", time.Minute, false},
		{"1h", time.Hour, false},
		{"1ns", time.Nanosecond, false},
		{"500ms", 500 * time.Millisecond, false},
		{"-500ms", -500 * time.Millisecond, false},
		{"1000ms", time.Second, false},
		{"1000000ns", time.Millisecond, false},
		{"1m1s", time.Minute + time.Second, false},
		{"9999999999ns", 9999999999 * time.Nanosecond, false},
		{"1.5ms", 1500 * time.Microsecond, false},
		{"1000ms", time.Second, false},
		{"-1000ms", -time.Second, false},
		{"1h1m1s", time.Hour + time.Minute + time.Second, false},
		{"3600s", time.Hour, false},
		{".1s", 100 * time.Millisecond, false},

		{"invalid", 0, true},
		{"-.g", 0, true},
		{"", 0, true},
		{"+", 0, true},
		{"1x", 0, true},
		{"1.", 0, true},
		{"1..1", 0, true},
		{"1a", 0, true},
		{"1.5", 0, true},
		{"s", 0, true},
		{"1..5s", 0, true},
		{"1h2m3.4s", time.Hour + 2*time.Minute + 3*time.Second + 400*time.Millisecond, false},
		{"1.0h", time.Hour, false},
		{"1.5h", time.Hour + 30*time.Minute, false},

		{"-1h", -time.Hour, false},
		{"0.5s", 500 * time.Millisecond, false},
		{"0.01s", 10 * time.Millisecond, false},
		{"99999ns", 99999 * time.Nanosecond, false},
		{"-9223372036854775808ns", 0, true},
		{"9223372038s", 0, true},
		{"9223372036854775807ns", time.Duration(9223372036854775807), false},
		{"91223372036854775807ns", 0 * time.Second, true},
		{"-91223372036854775808ns", 0, true},
		{"1.9122337203685477586776807s", 1*time.Second + 912233720*time.Nanosecond, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseDuration(%q) error = %v, wantErr %v", tt.input, err, tt.err)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseDuration(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

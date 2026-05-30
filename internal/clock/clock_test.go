package clock

import (
	"testing"
	"time"
)

func TestNewSnapshot(t *testing.T) {
	tests := []struct {
		name      string
		t         time.Time
		wantH12   int
		wantPM    bool
		wantBlink bool
	}{
		{"midnight", time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), 12, false, true},
		{"noon", time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC), 12, true, true},
		{"one pm", time.Date(2026, 1, 2, 13, 30, 1, 0, time.UTC), 1, true, false},
		{"eleven pm", time.Date(2026, 1, 2, 23, 59, 59, 0, time.UTC), 11, true, false},
		{"nine am even sec", time.Date(2026, 1, 2, 9, 8, 4, 0, time.UTC), 9, false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSnapshot(tc.t)
			if s.Hour12 != tc.wantH12 {
				t.Errorf("Hour12 = %d, want %d", s.Hour12, tc.wantH12)
			}
			if s.IsPM != tc.wantPM {
				t.Errorf("IsPM = %v, want %v", s.IsPM, tc.wantPM)
			}
			if s.BlinkOn != tc.wantBlink {
				t.Errorf("BlinkOn = %v, want %v", s.BlinkOn, tc.wantBlink)
			}
			if s.Hour24 != tc.t.Hour() {
				t.Errorf("Hour24 = %d, want %d", s.Hour24, tc.t.Hour())
			}
		})
	}
}

func TestParseMode(t *testing.T) {
	if m, err := ParseMode("analog"); err != nil || m != ModeAnalog {
		t.Errorf("ParseMode(analog) = %v, %v", m, err)
	}
	if _, err := ParseMode("bogus"); err == nil {
		t.Error("ParseMode(bogus) want error")
	}
}

func TestParseGranularity(t *testing.T) {
	if g, err := ParseGranularity("minutes"); err != nil || g != GranMinutes {
		t.Errorf("ParseGranularity(minutes) = %v, %v", g, err)
	}
	if _, err := ParseGranularity("hourly"); err == nil {
		t.Error("ParseGranularity(hourly) want error")
	}
}

// Package clock holds the framework-independent time domain shared by config,
// render, and ui. It imports none of them, which keeps those packages free of
// import cycles (enums live here,採用 #6).
package clock

import (
	"fmt"
	"time"
)

// Mode selects which renderer draws the clock.
type Mode uint8

const (
	ModeDigital Mode = iota
	ModeAnalog
)

func (m Mode) String() string {
	switch m {
	case ModeDigital:
		return "digital"
	case ModeAnalog:
		return "analog"
	default:
		return "unknown"
	}
}

// ParseMode maps a config string to a Mode.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "digital":
		return ModeDigital, nil
	case "analog":
		return ModeAnalog, nil
	default:
		return 0, fmt.Errorf("unknown mode %q (want \"digital\" or \"analog\")", s)
	}
}

// Granularity is how often the clock updates.
type Granularity uint8

const (
	GranSeconds Granularity = iota
	GranMinutes
)

func (g Granularity) String() string {
	switch g {
	case GranSeconds:
		return "seconds"
	case GranMinutes:
		return "minutes"
	default:
		return "unknown"
	}
}

// ParseGranularity maps a config string to a Granularity.
func ParseGranularity(s string) (Granularity, error) {
	switch s {
	case "seconds":
		return GranSeconds, nil
	case "minutes":
		return GranMinutes, nil
	default:
		return 0, fmt.Errorf("unknown granularity %q (want \"seconds\" or \"minutes\")", s)
	}
}

// TimeSnapshot is the pure data the View needs about "now". Deriving it here
// keeps View free of time.Now(), so rendering stays a pure function of state.
type TimeSnapshot struct {
	T       time.Time // source instant (local time)
	Hour24  int
	Hour12  int
	Minute  int
	Second  int
	IsPM    bool
	BlinkOn bool // colon-blink phase: true on even seconds
}

// NewSnapshot derives display fields from a source instant (local time).
func NewSnapshot(t time.Time) TimeSnapshot {
	h24 := t.Hour()
	h12 := h24 % 12
	if h12 == 0 {
		h12 = 12
	}
	return TimeSnapshot{
		T:       t,
		Hour24:  h24,
		Hour12:  h12,
		Minute:  t.Minute(),
		Second:  t.Second(),
		IsPM:    h24 >= 12,
		BlinkOn: t.Second()%2 == 0,
	}
}

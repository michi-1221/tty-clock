package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michi-1221/tty-clock/internal/clock"
)

// tickMsg carries the instant of a tick plus the generation that scheduled it.
// bubbletea has no way to cancel an in-flight tea.Every, so after a
// granularity change we self-dedupe: stale generations are dropped (★2).
type tickMsg struct {
	t   time.Time
	gen uint64
}

func interval(g clock.Granularity) time.Duration {
	if g == clock.GranMinutes {
		return time.Minute
	}
	return time.Second
}

// tickCmd schedules the next tick. tea.Every self-aligns to the wall-clock
// boundary (:00) so the display never drifts and DST/TZ self-correct.
func tickCmd(g clock.Granularity, gen uint64) tea.Cmd {
	return tea.Every(interval(g), func(t time.Time) tea.Msg {
		return tickMsg{t: t, gen: gen}
	})
}

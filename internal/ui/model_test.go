package ui

import (
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/michi-1221/tui-clock/internal/caps"
	"github.com/michi-1221/tui-clock/internal/config"
)

func testCaps() caps.Capabilities {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(termenv.Ascii)
	return caps.New(r, false)
}

func newModel(t *testing.T) Model {
	t.Helper()
	return New(config.DefaultConfig(), "", testCaps())
}

func runeKey(r rune) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}} }

func TestNewSeedsNow(t *testing.T) {
	m := newModel(t)
	if m.now.IsZero() {
		t.Error("New() must seed now (★1) so the first frame is correct")
	}
}

func TestTickReArms(t *testing.T) {
	m := newModel(t)
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	nm, cmd := m.Update(tickMsg{t: at, gen: m.tickGen})
	got := nm.(Model)
	if !got.now.Equal(at) {
		t.Errorf("now = %v, want %v", got.now, at)
	}
	if cmd == nil {
		t.Error("a current-generation tick MUST return a re-arm command (★2)")
	}
}

func TestStaleTickDropped(t *testing.T) {
	m := newModel(t)
	orig := m.now
	nm, cmd := m.Update(tickMsg{t: time.Now().Add(time.Hour), gen: m.tickGen + 99})
	got := nm.(Model)
	if !got.now.Equal(orig) {
		t.Error("stale-generation tick must not change now")
	}
	if cmd != nil {
		t.Error("stale-generation tick must not re-arm (would multiply loops)")
	}
}

func TestToggleSeconds(t *testing.T) {
	m := newModel(t)
	before := m.fmtOpts.ShowSeconds
	nm, _ := m.Update(runeKey('s'))
	if nm.(Model).fmtOpts.ShowSeconds == before {
		t.Error("'s' should toggle ShowSeconds")
	}
}

func TestCycleTheme(t *testing.T) {
	m := newModel(t)
	before := m.theme.Name
	nm, _ := m.Update(runeKey('t'))
	if nm.(Model).theme.Name == before {
		t.Errorf("'t' should change theme from %q", before)
	}
}

func TestToggleHelp(t *testing.T) {
	m := newModel(t)
	if !m.showHelp {
		t.Fatal("help line should be visible by default")
	}
	nm, _ := m.Update(runeKey('?'))
	if nm.(Model).showHelp {
		t.Error("'?' should hide the help line when visible")
	}
	nm2, _ := nm.(Model).Update(runeKey('?'))
	if !nm2.(Model).showHelp {
		t.Error("'?' again should show the help line")
	}
}

func TestHiddenHelpRemovesFooter(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	withHelp := nm.(Model).View()
	hidden, _ := nm.(Model).Update(runeKey('?'))
	withoutHelp := hidden.(Model).View()
	if strings.Contains(withoutHelp, "seconds") || strings.Contains(withoutHelp, "theme") {
		t.Error("hiding help should remove the seconds/theme labels from the footer")
	}
	if lipgloss.Height(withoutHelp) < lipgloss.Height(withHelp) {
		// sanity: hiding the footer shouldn't add rows
		t.Logf("with help %d rows, without %d", lipgloss.Height(withHelp), lipgloss.Height(withoutHelp))
	}
}

func TestQuit(t *testing.T) {
	m := newModel(t)
	nm, cmd := m.Update(runeKey('q'))
	if !nm.(Model).quitting {
		t.Error("'q' should set quitting")
	}
	if cmd == nil {
		t.Error("'q' should return tea.Quit")
	}
	if msg := cmd(); msg == nil {
		t.Error("quit command should yield a message")
	}
}

func TestModeToggleIsNoOp(t *testing.T) {
	m := newModel(t)
	nm, cmd := m.Update(runeKey('m'))
	if cmd != nil {
		t.Error("phase-1 'm' should be a no-op (★8)")
	}
	if nm.(Model).renderer.Name() != "digital" {
		t.Error("'m' must not switch away from digital in phase 1")
	}
}

func TestViewWithoutSizeDoesNotPanic(t *testing.T) {
	m := newModel(t)
	// No WindowSizeMsg yet: width/height are 0; must fall back, not break (★6).
	out := m.View()
	if strings.TrimSpace(out) == "" {
		t.Error("View with no size should still render (default size fallback)")
	}
}

func TestViewTinyTerminal(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(tea.WindowSizeMsg{Width: 3, Height: 2})
	out := nm.(Model).View()
	if strings.TrimSpace(out) == "" {
		t.Error("tiny terminal should still render something (no panic, no blank)")
	}
}

func TestScaleFor(t *testing.T) {
	cases := []struct{ w, h, want int }{
		{54, 10, 1},  // reference → ×1
		{108, 20, 2}, // both 2× → ×2
		{162, 30, 3}, // both 3× → ×3
		{110, 15, 1}, // height short → height-limited to ×1
		{40, 8, 1},   // below reference → floor at 1
		{53, 9, 1},
	}
	for _, c := range cases {
		if got := scaleFor(c.w, c.h); got != c.want {
			t.Errorf("scaleFor(%d,%d) = %d, want %d", c.w, c.h, got, c.want)
		}
	}
}

func TestLargeWindowScalesBody(t *testing.T) {
	m := newModel(t)
	small := m.body(60, 11, scaleFor(60, 12))
	big := m.body(130, 25, scaleFor(130, 26))
	if lipgloss.Height(big) <= lipgloss.Height(small) {
		t.Errorf("body height small=%d big=%d: a larger window should scale the clock up",
			lipgloss.Height(small), lipgloss.Height(big))
	}
}

func TestQuittingViewEmpty(t *testing.T) {
	m := newModel(t)
	nm, _ := m.Update(runeKey('q'))
	if nm.(Model).View() != "" {
		t.Error("quitting View should be empty to avoid a final flash")
	}
}

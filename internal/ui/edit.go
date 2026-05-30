package ui

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// editorFinishedMsg is delivered after the external editor (or OS opener) exits,
// so Update can reload the (possibly edited) config and surface any launch error.
type editorFinishedMsg struct{ err error }

var (
	errNoConfigPath = errors.New("no config file to open yet (running on built-in defaults — press r first)")
	errNoOpener     = errors.New("don't know how to open files on this OS; set $EDITOR")
)

// editConfigCmd opens path for editing. It prefers $VISUAL then $EDITOR: a
// terminal editor is run via tea.ExecProcess, which suspends the clock, hands
// over the TTY, and resumes on exit (Update then reloads the file). With neither
// set it falls back to the OS default opener (open / xdg-open / start), launched
// in the background without suspending — the user reloads later with 'r'.
func editConfigCmd(path string) tea.Cmd {
	if path == "" {
		return func() tea.Msg { return editorFinishedMsg{errNoConfigPath} }
	}
	if editor := firstEnv("VISUAL", "EDITOR"); editor != "" {
		fields := strings.Fields(editor) // honor flags, e.g. "code -w" / "emacsclient -nw"
		c := exec.Command(fields[0], append(fields[1:], path)...)
		return tea.ExecProcess(c, func(err error) tea.Msg { return editorFinishedMsg{err} })
	}
	c := openerCommand(path)
	if c == nil {
		return func() tea.Msg { return editorFinishedMsg{errNoOpener} }
	}
	return func() tea.Msg { return editorFinishedMsg{c.Start()} } // detach; don't suspend the clock
}

// firstEnv returns the first non-empty environment variable among keys.
func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

// openerCommand returns the OS default "open this file" command, or nil if the
// platform is unknown.
func openerCommand(path string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path)
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path)
	case "linux", "freebsd", "openbsd", "netbsd":
		return exec.Command("xdg-open", path)
	default:
		return nil
	}
}

package ui

import (
	"runtime"
	"testing"
)

func TestFirstEnvPrefersVisual(t *testing.T) {
	t.Setenv("VISUAL", "vim")
	t.Setenv("EDITOR", "nano")
	if got := firstEnv("VISUAL", "EDITOR"); got != "vim" {
		t.Errorf("firstEnv = %q, want vim", got)
	}
}

func TestFirstEnvFallsBackToEditor(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "nano")
	if got := firstEnv("VISUAL", "EDITOR"); got != "nano" {
		t.Errorf("firstEnv = %q, want nano", got)
	}
}

func TestOpenerCommandForCurrentOS(t *testing.T) {
	c := openerCommand("/tmp/x.json")
	want := map[string]string{
		"darwin": "open", "windows": "cmd",
		"linux": "xdg-open", "freebsd": "xdg-open", "openbsd": "xdg-open", "netbsd": "xdg-open",
	}[runtime.GOOS]
	if want == "" {
		t.Skipf("no opener expected on %s", runtime.GOOS)
	}
	if c == nil || c.Args[0] != want {
		t.Fatalf("opener on %s = %v, want %q", runtime.GOOS, c, want)
	}
	// The file path must be among the arguments.
	found := false
	for _, a := range c.Args {
		if a == "/tmp/x.json" {
			found = true
		}
	}
	if !found {
		t.Errorf("opener args %v missing the path", c.Args)
	}
}

func TestEditConfigCmdEmptyPathErrors(t *testing.T) {
	// Safe to invoke: the empty-path branch returns a message, never launches.
	msg, ok := editConfigCmd("")().(editorFinishedMsg)
	if !ok {
		t.Fatalf("want editorFinishedMsg, got %T", msg)
	}
	if msg.err == nil {
		t.Error("empty path should yield an error message")
	}
}

func TestEditConfigCmdWithEditorIsNonNil(t *testing.T) {
	// With $EDITOR set we get an ExecProcess command; we only assert it exists
	// (running it would launch an editor and grab the TTY).
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "true") // a real, harmless binary
	if editConfigCmd("/tmp/x.json") == nil {
		t.Error("editConfigCmd should return a command when $EDITOR is set")
	}
}

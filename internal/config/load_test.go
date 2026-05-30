package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cfgPath(name string) string {
	return filepath.Join("..", "..", "testdata", "configs", name)
}

func TestLoadValid(t *testing.T) {
	cfg, err := Load(cfgPath("valid.json"))
	if err != nil {
		t.Fatalf("Load(valid) error: %v", err)
	}
	if cfg.Theme != "dracula" {
		t.Errorf("Theme = %q, want dracula", cfg.Theme)
	}
	if cfg.Format.Hour24 {
		t.Error("Hour24 = true, want false (from file)")
	}
	if cfg.Format.BlinkColon != true {
		t.Error("BlinkColon = false, want true (from file)")
	}
	// Keys not present in the file keep their defaults.
	if cfg.Format.DateFormat != "Mon 2006-01-02" {
		t.Errorf("DateFormat = %q, want default kept", cfg.Format.DateFormat)
	}
	if !cfg.Format.ShowNumbers {
		t.Error("ShowNumbers should keep default true")
	}
	if cfg.Format.Font != "block" {
		t.Errorf("Font = %q, want default block", cfg.Format.Font)
	}
}

func TestLoadPartialKeepsDefaults(t *testing.T) {
	cfg, err := Load(cfgPath("partial.json"))
	if err != nil {
		t.Fatalf("Load(partial) error: %v", err)
	}
	if cfg.Theme != "nord" {
		t.Errorf("Theme = %q, want nord", cfg.Theme)
	}
	if !cfg.Format.Hour24 || !cfg.Format.ShowSeconds || cfg.Mode != "digital" {
		t.Errorf("defaults not kept: %+v", cfg)
	}
}

func TestLoadBadEnum(t *testing.T) {
	_, err := Load(cfgPath("badenum.json"))
	if err == nil || !strings.Contains(err.Error(), "granularity") {
		t.Fatalf("want granularity error, got: %v", err)
	}
}

func TestLoadUnknownKey(t *testing.T) {
	_, err := Load(cfgPath("unknownkey.json"))
	if err == nil || !strings.Contains(err.Error(), "wat") {
		t.Fatalf("want unknown-field error naming 'wat', got: %v", err)
	}
}

func TestLoadSyntaxError(t *testing.T) {
	_, err := Load(cfgPath("syntaxerr.json"))
	if err == nil || !strings.Contains(err.Error(), "line") {
		t.Fatalf("want syntax error with line:col, got: %v", err)
	}
}

func TestLoadTypeError(t *testing.T) {
	_, err := Load(cfgPath("typeerr.json"))
	if err == nil || !strings.Contains(err.Error(), "type") {
		t.Fatalf("want type error, got: %v", err)
	}
}

func TestResolveScaffoldsOnFirstRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := Resolve("")
	if err != nil {
		t.Fatalf("Resolve(\"\") error: %v", err)
	}
	want := filepath.Join(home, ".tty-clock", "config.json")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("scaffolded file missing: %v", err)
	}
	// The scaffolded file must load and equal the built-in defaults.
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load(scaffolded) error: %v", err)
	}
	if cfg != DefaultConfig() {
		t.Errorf("scaffolded config != defaults:\n got %+v\nwant %+v", cfg, DefaultConfig())
	}
}

func TestResolveDoesNotOverwriteExisting(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".tty-clock", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const custom = `{"theme":"nord"}`
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := Resolve("")
	if err != nil || got != path {
		t.Fatalf("Resolve(\"\") = %q, %v; want %q", got, err, path)
	}
	if b, _ := os.ReadFile(path); string(b) != custom {
		t.Errorf("Resolve overwrote an existing config: %q", string(b))
	}
}

func TestResolveExplicitMissing(t *testing.T) {
	if _, err := Resolve(cfgPath("does-not-exist.json")); err == nil {
		t.Error("explicit missing config should be fatal")
	}
}

func TestResolveExplicitExisting(t *testing.T) {
	p, err := Resolve(cfgPath("valid.json"))
	if err != nil || p == "" {
		t.Fatalf("Resolve(existing) = %q, %v", p, err)
	}
}

func TestLoadEmptyPathUsesDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Theme != DefaultConfig().Theme {
		t.Error("empty path should yield defaults")
	}
}

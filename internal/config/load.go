package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".tty-clock" // created under the user's home directory
	configFileName = "config.json"
)

// HomeConfigPath returns ~/.tty-clock/config.json.
func HomeConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName, configFileName), nil
}

// Resolve decides which config file to load:
//   - explicit != "": that file MUST exist (a missing file is fatal).
//   - explicit == "": use ~/.tty-clock/config.json, creating it with the default
//     config on first run, then loading it.
//
// Scaffolding is best-effort: if the home dir is unknown or the file can't be
// written, Resolve returns "" so the caller runs with built-in defaults rather
// than failing to start.
func Resolve(explicit string) (string, error) {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return "", fmt.Errorf("config file %q: %w", explicit, err)
		}
		return explicit, nil
	}
	path, err := HomeConfigPath()
	if err != nil {
		return "", nil // home unknown → built-in defaults
	}
	switch _, statErr := os.Stat(path); {
	case statErr == nil:
		return path, nil // already exists → use it (never overwrite)
	case !os.IsNotExist(statErr):
		return "", fmt.Errorf("config file %q: %w", path, statErr)
	}
	if scaffold(path) != nil {
		return "", nil // couldn't write → built-in defaults
	}
	return path, nil
}

// scaffold writes the default config to path (creating ~/.tty-clock/), so a
// first-time user gets an editable starting point.
func scaffold(path string) error {
	return Save(path, DefaultConfig())
}

// Load reads and validates the config at path. When path is "", it returns the
// defaults unchanged.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config %q: %w", path, err)
	}
	if err := decodeInto(&cfg, data); err != nil {
		return cfg, fmt.Errorf("%s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

// decodeInto overlays JSON onto an already-defaulted config (so missing keys,
// including nested Format keys, keep their defaults) and turns the stdlib's
// terse errors into line:col messages.
func decodeInto(cfg *Config, data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(cfg); err != nil {
		return friendlyJSONError(data, err)
	}
	// Reject trailing content after the first JSON value.
	if err := dec.Decode(new(json.RawMessage)); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing data after the JSON object")
		}
		return friendlyJSONError(data, err)
	}
	return nil
}

// friendlyJSONError dispatches over the three error shapes the decoder can
// produce: syntax errors, type mismatches, and unknown fields (★8).
func friendlyJSONError(data []byte, err error) error {
	var syn *json.SyntaxError
	if errors.As(err, &syn) {
		line, col := offsetToLineCol(data, syn.Offset)
		return fmt.Errorf("JSON syntax error at line %d:%d: %s", line, col, syn.Error())
	}
	var typ *json.UnmarshalTypeError
	if errors.As(err, &typ) {
		line, col := offsetToLineCol(data, typ.Offset)
		return fmt.Errorf("JSON type error at line %d:%d: field %q expects %s but got %s",
			line, col, typ.Field, typ.Type, typ.Value)
	}
	// DisallowUnknownFields and the rest already name the offending field.
	return err
}

// offsetToLineCol converts a byte offset (1-based line/col) for error messages.
func offsetToLineCol(data []byte, off int64) (line, col int) {
	line, col = 1, 1
	for i := int64(0); i < off && i < int64(len(data)); i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

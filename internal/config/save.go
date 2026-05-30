package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Save writes cfg to path as pretty JSON, atomically (temp file + rename) so a
// crash mid-write can't corrupt the config. It creates the parent directory if
// needed.
func Save(path string, cfg Config) error {
	if path == "" {
		return errors.New("no config path to save to")
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

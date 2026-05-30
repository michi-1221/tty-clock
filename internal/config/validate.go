package config

import (
	"fmt"

	"github.com/michi-1221/tty-clock/internal/clock"
)

// Validate checks enum fields, returning a helpful error that lists the
// allowed values. Theme existence is validated later by theme.Resolve.
func (c Config) Validate() error {
	if _, err := clock.ParseMode(c.Mode); err != nil {
		return err
	}
	if _, err := clock.ParseGranularity(c.Granularity); err != nil {
		return err
	}
	switch c.Format.Font {
	case "block", "ascii", "seg7":
	default:
		return fmt.Errorf("unknown font %q (want \"block\", \"ascii\" or \"seg7\")", c.Format.Font)
	}
	return nil
}

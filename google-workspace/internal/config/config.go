package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CalendarMode controls the level of Calendar API access.
type CalendarMode string

const (
	CalendarOff       CalendarMode = "off"
	CalendarReadOnly  CalendarMode = "readonly"
	CalendarReadWrite CalendarMode = "readwrite"
)

// Config holds the scope configuration for the skill.
type Config struct {
	Gmail    bool         `json:"gmail"`
	Calendar CalendarMode `json:"calendar"`
	Contacts bool         `json:"contacts"`
	Drive    bool         `json:"drive"`
}

// DefaultConfig returns the default (most restrictive) configuration.
func DefaultConfig() Config {
	return Config{
		Gmail:    true,
		Calendar: CalendarReadOnly,
		Contacts: true,
		Drive:    true,
	}
}

// Validate checks that the config values are within expected bounds.
func (c Config) Validate() error {
	switch c.Calendar {
	case CalendarOff, CalendarReadOnly, CalendarReadWrite:
		// valid
	default:
		return fmt.Errorf("invalid calendar mode %q: must be off, readonly, or readwrite", c.Calendar)
	}
	return nil
}

// Load reads config from the given directory. If the file does not exist,
// it returns the default config without error.
func Load(dir string) (Config, error) {
	path := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Save writes the config to the given directory, creating it if needed.
func Save(dir string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	data = append(data, '\n')
	path := filepath.Join(dir, "config.json")

	return os.WriteFile(path, data, 0o600)
}

// OAuthScopes returns the Google OAuth scopes needed for the current config.
func (c Config) OAuthScopes() []string {
	var scopes []string

	if c.Gmail {
		scopes = append(scopes, "https://www.googleapis.com/auth/gmail.readonly")
	}

	switch c.Calendar {
	case CalendarReadOnly:
		scopes = append(scopes, "https://www.googleapis.com/auth/calendar.readonly")
	case CalendarReadWrite:
		scopes = append(scopes, "https://www.googleapis.com/auth/calendar.events")
	}

	if c.Contacts {
		scopes = append(scopes, "https://www.googleapis.com/auth/contacts.readonly")
	}

	if c.Drive {
		scopes = append(scopes, "https://www.googleapis.com/auth/drive.readonly")
	}

	return scopes
}

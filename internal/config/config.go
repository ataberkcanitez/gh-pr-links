// Package config loads and persists the user's gh-pr-links preferences.
//
// The configuration lives in a JSON file in the user's home directory. Its
// on-disk format is part of the public contract of the extension: existing
// users already have this file, so field names and value types must not change.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileName is the name of the configuration file inside the home directory.
const FileName = "gh-pr-links-configuration.json"

// BypassFilter is the organization value that disables organization filtering.
const BypassFilter = "-"

// Config holds the user's preferences.
//
// UseEmoji is stored as a string ("true"/"false") rather than a bool to stay
// compatible with configuration files written by earlier releases.
type Config struct {
	Style        string `json:"style"`
	UseEmoji     string `json:"useEmoji"`
	Organization string `json:"organization"`
}

// Default returns the configuration used when no file exists yet.
func Default() Config {
	return Config{
		Style:        "StyleRounded",
		UseEmoji:     "true",
		Organization: BypassFilter,
	}
}

// EmojiEnabled reports whether emojis should be shown in the output.
func (c Config) EmojiEnabled() bool {
	return c.UseEmoji == "true"
}

// Organizations returns the organizations to show pull requests from, or nil
// when filtering is disabled. Multiple organizations are separated by "|".
func (c Config) Organizations() []string {
	if c.Organization == BypassFilter {
		return nil
	}
	return strings.Split(c.Organization, "|")
}

// Store reads and writes a Config at a fixed path.
type Store struct {
	path string
}

// NewStore returns a Store backed by the file at path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// DefaultPath returns the location of the configuration file in the user's
// home directory.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locate home directory: %w", err)
	}
	return filepath.Join(home, FileName), nil
}

// Path returns the file the store reads from and writes to.
func (s *Store) Path() string {
	return s.path
}

// Load reads the configuration, creating the file with defaults if it does
// not exist yet.
func (s *Store) Load() (Config, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		cfg := Default()
		if err := s.Save(cfg); err != nil {
			return Config{}, err
		}
		return cfg, nil
	}
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	decodeErr := json.Unmarshal(data, &cfg)

	// Files written before organization filtering existed have no
	// organization; migrate them to "no filter". This intentionally runs even
	// when decoding failed, which rewrites a corrupt file so the next run
	// recovers on its own.
	if cfg.Organization == "" {
		cfg.Organization = BypassFilter
		_ = s.Save(cfg)
	}
	if decodeErr != nil {
		return Config{}, decodeErr
	}
	return cfg, nil
}

// Save writes cfg to disk, replacing any existing configuration.
func (s *Store) Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

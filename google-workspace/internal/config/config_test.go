package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.Gmail {
		t.Error("expected gmail enabled by default")
	}
	if cfg.Calendar != CalendarReadOnly {
		t.Errorf("expected calendar readonly by default, got %s", cfg.Calendar)
	}
	if !cfg.Contacts {
		t.Error("expected contacts enabled by default")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		mode    CalendarMode
		wantErr bool
	}{
		{"off", CalendarOff, false},
		{"readonly", CalendarReadOnly, false},
		{"readwrite", CalendarReadWrite, false},
		{"invalid", CalendarMode("bogus"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Gmail: true, Calendar: tt.mode, Contacts: true}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOAuthScopes(t *testing.T) {
	tests := []struct {
		name   string
		cfg    Config
		expect int
	}{
		{"all enabled readonly", Config{Gmail: true, Calendar: CalendarReadOnly, Contacts: true}, 3},
		{"all enabled readwrite", Config{Gmail: true, Calendar: CalendarReadWrite, Contacts: true}, 3},
		{"gmail only", Config{Gmail: true, Calendar: CalendarOff, Contacts: false}, 1},
		{"nothing", Config{Gmail: false, Calendar: CalendarOff, Contacts: false}, 0},
		{"calendar only", Config{Gmail: false, Calendar: CalendarReadOnly, Contacts: false}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes := tt.cfg.OAuthScopes()
			if len(scopes) != tt.expect {
				t.Errorf("expected %d scopes, got %d: %v", tt.expect, len(scopes), scopes)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	cfg := Config{
		Gmail:    true,
		Calendar: CalendarReadWrite,
		Contacts: false,
	}

	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Gmail != cfg.Gmail || loaded.Calendar != cfg.Calendar || loaded.Contacts != cfg.Contacts {
		t.Errorf("loaded config does not match saved: got %+v, want %+v", loaded, cfg)
	}
}

func TestLoadMissing(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	def := DefaultConfig()
	if cfg != def {
		t.Errorf("expected default config for missing file, got %+v", cfg)
	}
}

func TestLoadInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"calendar":"bogus"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error loading invalid config")
	}
}

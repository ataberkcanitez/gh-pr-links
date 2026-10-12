package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewStore(filepath.Join(t.TempDir(), FileName))
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCreatesDefaultFile(t *testing.T) {
	store := newTestStore(t)

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg != Default() {
		t.Errorf("Load() = %+v, want %+v", cfg, Default())
	}

	// The on-disk format is shared with earlier releases and must not change.
	want := "{\n  \"style\": \"StyleRounded\",\n  \"useEmoji\": \"true\",\n  \"organization\": \"-\"\n}"
	if got := readFile(t, store.Path()); got != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}

func TestLoadExistingFile(t *testing.T) {
	store := newTestStore(t)
	writeFile(t, store.Path(), `{"style":"StyleMarkdown","useEmoji":"false","organization":"acme|globex"}`)

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Config{Style: "StyleMarkdown", UseEmoji: "false", Organization: "acme|globex"}
	if cfg != want {
		t.Errorf("Load() = %+v, want %+v", cfg, want)
	}
}

func TestLoadMigratesMissingOrganization(t *testing.T) {
	store := newTestStore(t)
	writeFile(t, store.Path(), `{"style":"StyleUnicode","useEmoji":"true"}`)

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Organization != BypassFilter {
		t.Errorf("Organization = %q, want %q", cfg.Organization, BypassFilter)
	}

	reloaded, err := store.Load()
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if reloaded != cfg {
		t.Errorf("migrated config was not persisted: got %+v, want %+v", reloaded, cfg)
	}
}

func TestLoadCorruptFileRecoversOnNextRun(t *testing.T) {
	store := newTestStore(t)
	writeFile(t, store.Path(), `{not json`)

	if _, err := store.Load(); err == nil {
		t.Fatal("Load succeeded on a corrupt file, want error")
	}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if want := (Config{Organization: BypassFilter}); cfg != want {
		t.Errorf("Load() = %+v, want %+v", cfg, want)
	}
}

func TestLoadUnwritableLocation(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "missing-dir", FileName))
	if _, err := store.Load(); err == nil {
		t.Fatal("Load succeeded, want error when the default file cannot be created")
	}
}

func TestSaveRoundTrip(t *testing.T) {
	store := newTestStore(t)
	want := Config{Style: "StyleCompact", UseEmoji: "false", Organization: "acme"}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestEmojiEnabled(t *testing.T) {
	for value, want := range map[string]bool{"true": true, "false": false, "": false, "TRUE": false} {
		if got := (Config{UseEmoji: value}).EmojiEnabled(); got != want {
			t.Errorf("EmojiEnabled(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestOrganizations(t *testing.T) {
	tests := map[string][]string{
		BypassFilter:   nil,
		"acme":         {"acme"},
		"acme|globex":  {"acme", "globex"},
		"acme||globex": {"acme", "", "globex"},
	}
	for value, want := range tests {
		if got := (Config{Organization: value}).Organizations(); !reflect.DeepEqual(got, want) {
			t.Errorf("Organizations(%q) = %q, want %q", value, got, want)
		}
	}
}

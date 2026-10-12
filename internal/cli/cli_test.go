package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ataberkcanitez/gh-pr-links/internal/config"
	"github.com/ataberkcanitez/gh-pr-links/internal/github"
)

type fakeStore struct {
	cfg     config.Config
	loadErr error
	saveErr error
	saves   int
}

func (s *fakeStore) Load() (config.Config, error) { return s.cfg, s.loadErr }

func (s *fakeStore) Save(cfg config.Config) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.cfg = cfg
	s.saves++
	return nil
}

type fakeSource struct {
	prs    []github.PullRequest
	err    error
	called bool
}

func (s *fakeSource) ReviewRequested() ([]github.PullRequest, error) {
	s.called = true
	return s.prs, s.err
}

type harness struct {
	app    *App
	store  *fakeStore
	source *fakeSource
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func newHarness() *harness {
	h := &harness{
		store:  &fakeStore{cfg: config.Default()},
		source: &fakeSource{},
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}
	h.app = &App{Name: "gh-pr-links", Stdout: h.stdout, Stderr: h.stderr, Config: h.store, PRs: h.source}
	return h
}

func pr(repo, title string) github.PullRequest {
	return github.PullRequest{
		Author:     github.Author{Login: "alice"},
		Repository: github.Repository{NameWithOwner: repo},
		Title:      title,
		URL:        "https://github.com/" + repo + "/pull/1",
	}
}

func TestHelp(t *testing.T) {
	for _, flag := range []string{"-h", "--help", "-help"} {
		t.Run(flag, func(t *testing.T) {
			h := newHarness()
			if code := h.app.Run([]string{flag}); code != ExitOK {
				t.Errorf("exit code = %d, want %d", code, ExitOK)
			}
			want := "pr-links is a command line tool to show your open and review requested pull requests.\n" +
				"Usage: gh pr-links [options]\n" +
				"Options:\n" +
				"  -h, --help\t\t\tShows help message\n" +
				"  -s, --style <string>\t\tSets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic\n" +
				"  -u, --use-emoji <bool>\tUse emoji in output. Possible values: true, false\n" +
				"  -o, --org <string>\tFilter for organization, use '-' for skip the filtering, or <org_name> to filter based on organization\n" +
				"\nConfiguration: gh pr-links --style=StyleCompactLite --use-emoji=true\n" +
				"Usage: gh pr-links\n"
			if got := h.stdout.String(); got != want {
				t.Errorf("help output:\n%s\nwant:\n%s", got, want)
			}
			if h.source.called {
				t.Error("help should not fetch pull requests")
			}
		})
	}
}

func TestHelpTakesPrecedenceOverSettings(t *testing.T) {
	h := newHarness()
	h.app.Run([]string{"--help", "--style=StyleMarkdown", "--org=acme"})
	if h.store.saves != 0 {
		t.Errorf("saved config %d times, want 0", h.store.saves)
	}
}

func TestUnknownFlag(t *testing.T) {
	h := newHarness()
	if code := h.app.Run([]string{"--nope"}); code != ExitUsageError {
		t.Errorf("exit code = %d, want %d", code, ExitUsageError)
	}
	if !strings.Contains(h.stderr.String(), "flag provided but not defined: -nope") {
		t.Errorf("stderr = %q, want unknown flag error", h.stderr.String())
	}
	if !strings.Contains(h.stderr.String(), "Usage of gh-pr-links:") {
		t.Errorf("stderr = %q, want usage", h.stderr.String())
	}
}

func TestConfigLoadError(t *testing.T) {
	h := newHarness()
	h.store.loadErr = errors.New("disk on fire")
	if code := h.app.Run(nil); code != ExitError {
		t.Errorf("exit code = %d, want %d", code, ExitError)
	}
	if got, want := h.stdout.String(), "Error: failed to read config file: disk on fire\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestSetStyle(t *testing.T) {
	for _, flag := range []string{"--style=StyleMarkdown", "-s=StyleMarkdown"} {
		t.Run(flag, func(t *testing.T) {
			h := newHarness()
			h.source.prs = []github.PullRequest{pr("acme/api", "Fix")}
			if code := h.app.Run([]string{flag}); code != ExitOK {
				t.Errorf("exit code = %d, want %d", code, ExitOK)
			}
			if h.store.cfg.Style != "StyleMarkdown" {
				t.Errorf("saved style = %q, want StyleMarkdown", h.store.cfg.Style)
			}
			if !h.source.called {
				t.Error("setting a style should still list pull requests")
			}
		})
	}
}

func TestInvalidStyle(t *testing.T) {
	h := newHarness()
	if code := h.app.Run([]string{"--style=Fancy", "--use-emoji=false"}); code != ExitError {
		t.Errorf("exit code = %d, want %d", code, ExitError)
	}
	want := "Error: style value is invalid. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic\n"
	if got := h.stdout.String(); got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if h.store.saves != 0 {
		t.Error("an invalid style must not save anything")
	}
}

func TestSetEmoji(t *testing.T) {
	for _, value := range []string{"true", "false"} {
		t.Run(value, func(t *testing.T) {
			h := newHarness()
			h.store.cfg.UseEmoji = "unset"
			h.app.Run([]string{"-u", value})
			if h.store.cfg.UseEmoji != value {
				t.Errorf("saved useEmoji = %q, want %q", h.store.cfg.UseEmoji, value)
			}
		})
	}
}

func TestInvalidEmoji(t *testing.T) {
	h := newHarness()
	if code := h.app.Run([]string{"--use-emoji=yes"}); code != ExitError {
		t.Errorf("exit code = %d, want %d", code, ExitError)
	}
	if got, want := h.stdout.String(), "Error: use-emoji value is invalid. Possible values: true, false\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestSaveErrors(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"--style=StyleUnicode"}, "Error: failed to update config file: read-only\n"},
		{[]string{"--use-emoji=true"}, "Error: failed to update config file: read-only\n"},
		{[]string{"--org=acme"}, "Error: failed to update config: read-only\n"},
	}
	for _, tt := range tests {
		t.Run(tt.args[0], func(t *testing.T) {
			h := newHarness()
			h.store.saveErr = errors.New("read-only")
			if code := h.app.Run(tt.args); code != ExitError {
				t.Errorf("exit code = %d, want %d", code, ExitError)
			}
			if got := h.stdout.String(); got != tt.want {
				t.Errorf("stdout = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetOrganization(t *testing.T) {
	tests := []struct {
		args []string
		org  string
		want string
	}{
		{
			args: []string{"--org=acme|globex"},
			org:  "acme|globex",
			want: "Organization filter set to 'acme|globex'. Only pull requests from these organizations will be shown.\n",
		},
		{
			args: []string{"-o", "-"},
			org:  "-",
			want: "Organization filter removed. All pull requests will be shown\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.org, func(t *testing.T) {
			h := newHarness()
			if code := h.app.Run(tt.args); code != ExitOK {
				t.Errorf("exit code = %d, want %d", code, ExitOK)
			}
			if h.store.cfg.Organization != tt.org {
				t.Errorf("saved organization = %q, want %q", h.store.cfg.Organization, tt.org)
			}
			if got := h.stdout.String(); got != tt.want {
				t.Errorf("stdout = %q, want %q", got, tt.want)
			}
			if h.source.called {
				t.Error("setting the organization should not list pull requests")
			}
		})
	}
}

func TestSettingsCombineWithOrganization(t *testing.T) {
	h := newHarness()
	h.app.Run([]string{"--style=StyleUnicode", "--use-emoji=false", "--org=acme"})
	want := config.Config{Style: "StyleUnicode", UseEmoji: "false", Organization: "acme"}
	if h.store.cfg != want {
		t.Errorf("saved config = %+v, want %+v", h.store.cfg, want)
	}
}

func TestListFetchError(t *testing.T) {
	h := newHarness()
	h.source.err = errors.New("failed to execute gh command: exit status 1")
	if code := h.app.Run(nil); code != ExitError {
		t.Errorf("exit code = %d, want %d", code, ExitError)
	}
	if got, want := h.stdout.String(), "Error: failed to execute gh command: exit status 1\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestListEmpty(t *testing.T) {
	h := newHarness()
	if code := h.app.Run(nil); code != ExitOK {
		t.Errorf("exit code = %d, want %d", code, ExitOK)
	}
	if got, want := h.stdout.String(), "No PR to show.\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestListFiltersByOrganization(t *testing.T) {
	h := newHarness()
	h.store.cfg.Organization = "acme"
	h.source.prs = []github.PullRequest{pr("acme/api", "Keep me"), pr("globex/web", "Drop me")}

	if code := h.app.Run(nil); code != ExitOK {
		t.Errorf("exit code = %d, want %d", code, ExitOK)
	}
	out := h.stdout.String()
	if !strings.Contains(out, "acme/api") {
		t.Errorf("output is missing acme/api:\n%s", out)
	}
	if strings.Contains(out, "globex/web") {
		t.Errorf("output should not contain globex/web:\n%s", out)
	}
}

func TestListAllFilteredOut(t *testing.T) {
	h := newHarness()
	h.store.cfg.Organization = "acme"
	h.source.prs = []github.PullRequest{pr("globex/web", "Drop me")}

	h.app.Run(nil)
	if got, want := h.stdout.String(), "No PR to show.\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

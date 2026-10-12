package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ataberkcanitez/gh-pr-links/internal/github"
)

func TestRender(t *testing.T) {
	prs := []github.PullRequest{
		{
			Author:     github.Author{Login: "alice"},
			Repository: github.Repository{NameWithOwner: "acme/api"},
			Title:      "Fix bug",
			URL:        "https://github.com/acme/api/pull/1",
		},
		{
			Author:     github.Author{Login: "bob"},
			IsDraft:    true,
			Repository: github.Repository{NameWithOwner: "acme/web"},
			Title:      "WIP",
			URL:        "https://github.com/acme/web/pull/2",
		},
	}

	var buf bytes.Buffer
	if err := (Table{Style: "StyleMarkdown"}).Render(&buf, prs); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()

	for _, want := range []string{
		"REPO", "TITLE", "AUTHOR", "LINK",
		green("acme/api"), green("Fix"), green("alice"), blue("https://github.com/acme/api/pull/1"),
		gray("acme/web"), gray("WIP"), gray("bob"), gray("https://github.com/acme/web/pull/2"),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output is missing %q:\n%s", want, out)
		}
	}
	if !strings.HasSuffix(out, "\n") {
		t.Error("output should end with a newline")
	}
}

func TestTitle(t *testing.T) {
	tests := []struct {
		name     string
		useEmoji bool
		isDraft  bool
		title    string
		want     string
	}{
		{
			name:  "short title",
			title: "Fix bug",
			want:  "<Fix> <bug> ",
		},
		{
			name:     "ready emoji",
			useEmoji: true,
			title:    "Fix bug",
			want:     readyEmoji + " <Fix> <bug> ",
		},
		{
			name:     "draft emoji",
			useEmoji: true,
			isDraft:  true,
			title:    "Fix",
			want:     draftEmoji + " <Fix> ",
		},
		{
			name:  "collapses whitespace",
			title: "  Fix \t  bug  ",
			want:  "<Fix> <bug> ",
		},
		{
			name:  "wraps after 40 characters",
			title: "aaaaaaaaaa bbbbbbbbbb cccccccccc dddddddddd eeeee",
			want:  "<aaaaaaaaaa> <bbbbbbbbbb> <cccccccccc> \n<dddddddddd> <eeeee> ",
		},
		{
			name:  "long first word starts on a new line",
			title: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:  "\n<aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa> ",
		},
		{
			name:  "empty title",
			title: "",
			want:  "",
		},
	}

	mark := func(s string) string { return "<" + s + ">" }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Table{UseEmoji: tt.useEmoji}.title(tt.title, tt.isDraft, mark)
			if got != tt.want {
				t.Errorf("title(%q) = %q, want %q", tt.title, got, tt.want)
			}
		})
	}
}

func TestIsValidStyle(t *testing.T) {
	for _, name := range StyleNames() {
		if !IsValidStyle(name) {
			t.Errorf("IsValidStyle(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"", "rounded", "StyleUnknown"} {
		if IsValidStyle(name) {
			t.Errorf("IsValidStyle(%q) = true, want false", name)
		}
	}
}

func TestStyleNamesMatchesStyles(t *testing.T) {
	names := StyleNames()
	if len(names) != len(styles) {
		t.Fatalf("StyleNames has %d entries, styles map has %d", len(names), len(styles))
	}
	names[0] = "mutated"
	if StyleNames()[0] == "mutated" {
		t.Error("StyleNames returned the internal slice")
	}
}

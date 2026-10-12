package github

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestReviewRequested(t *testing.T) {
	var gotArgs []string
	client := NewClientWithExec(func(args ...string) (bytes.Buffer, bytes.Buffer, error) {
		gotArgs = args
		out := bytes.NewBufferString(`[
			{"author":{"login":"alice"},"isDraft":true,"repository":{"nameWithOwner":"acme/api"},"title":"Fix","url":"https://github.com/acme/api/pull/1"}
		]`)
		return *out, bytes.Buffer{}, nil
	})

	prs, err := client.ReviewRequested()
	if err != nil {
		t.Fatalf("ReviewRequested: %v", err)
	}

	wantArgs := []string{"search", "prs", "--review-requested", "@me", "--state", "open", "--json", "author,repository,isDraft,title,url"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("gh args = %q, want %q", gotArgs, wantArgs)
	}

	want := []PullRequest{{
		Author:     Author{Login: "alice"},
		IsDraft:    true,
		Repository: Repository{NameWithOwner: "acme/api"},
		Title:      "Fix",
		URL:        "https://github.com/acme/api/pull/1",
	}}
	if !reflect.DeepEqual(prs, want) {
		t.Errorf("prs = %+v, want %+v", prs, want)
	}
}

func TestReviewRequestedErrors(t *testing.T) {
	t.Run("gh fails", func(t *testing.T) {
		client := NewClientWithExec(func(...string) (bytes.Buffer, bytes.Buffer, error) {
			return bytes.Buffer{}, bytes.Buffer{}, errors.New("boom")
		})
		_, err := client.ReviewRequested()
		if err == nil || !strings.Contains(err.Error(), "failed to execute gh command: boom") {
			t.Errorf("err = %v, want gh execution error", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		client := NewClientWithExec(func(...string) (bytes.Buffer, bytes.Buffer, error) {
			return *bytes.NewBufferString("not json"), bytes.Buffer{}, nil
		})
		_, err := client.ReviewRequested()
		if err == nil || !strings.HasPrefix(err.Error(), "failed to unmarshal pull requests:") {
			t.Errorf("err = %v, want unmarshal error", err)
		}
	})
}

func TestOwner(t *testing.T) {
	tests := map[string]string{
		"acme/api": "acme",
		"acme":     "acme",
		"":         "",
	}
	for name, want := range tests {
		pr := PullRequest{Repository: Repository{NameWithOwner: name}}
		if got := pr.Owner(); got != want {
			t.Errorf("Owner(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestFilterByOrganizations(t *testing.T) {
	pr := func(repo string) PullRequest {
		return PullRequest{Repository: Repository{NameWithOwner: repo}}
	}
	prs := []PullRequest{pr("acme/api"), pr("globex/web"), pr("acme/cli"), pr("initech/tps")}

	tests := []struct {
		name string
		orgs []string
		want []PullRequest
	}{
		{"single org", []string{"acme"}, []PullRequest{pr("acme/api"), pr("acme/cli")}},
		{"multiple orgs keep order", []string{"initech", "acme"}, []PullRequest{pr("acme/api"), pr("acme/cli"), pr("initech/tps")}},
		{"case sensitive", []string{"ACME"}, nil},
		{"no match", []string{"umbrella"}, nil},
		{"no orgs", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByOrganizations(prs, tt.orgs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterByOrganizations(%q) = %v, want %v", tt.orgs, got, tt.want)
			}
		})
	}
}

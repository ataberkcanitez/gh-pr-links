// Package github fetches pull requests through the GitHub CLI.
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2"
)

// PullRequest is the subset of pull request fields the extension displays.
type PullRequest struct {
	Author     Author     `json:"author"`
	IsDraft    bool       `json:"isDraft"`
	Repository Repository `json:"repository"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
}

// Author is the user who opened a pull request.
type Author struct {
	Login string `json:"login"`
}

// Repository identifies the repository a pull request belongs to.
type Repository struct {
	NameWithOwner string `json:"nameWithOwner"`
}

// Owner returns the user or organization that owns the pull request's
// repository.
func (pr PullRequest) Owner() string {
	owner, _, _ := strings.Cut(pr.Repository.NameWithOwner, "/")
	return owner
}

// ExecFunc runs a gh command. Its signature matches gh.Exec.
type ExecFunc func(args ...string) (stdout, stderr bytes.Buffer, err error)

// Client queries pull requests via the gh CLI.
type Client struct {
	exec ExecFunc
}

// NewClient returns a Client that shells out to the installed gh CLI.
func NewClient() *Client {
	return NewClientWithExec(gh.Exec)
}

// NewClientWithExec returns a Client that runs gh commands through exec.
// It is mainly useful for tests.
func NewClientWithExec(exec ExecFunc) *Client {
	return &Client{exec: exec}
}

// ReviewRequested returns the open pull requests that request a review from
// the authenticated user.
func (c *Client) ReviewRequested() ([]PullRequest, error) {
	stdout, _, err := c.exec(
		"search", "prs",
		"--review-requested", "@me",
		"--state", "open",
		"--json", "author,repository,isDraft,title,url",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}

	var prs []PullRequest
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pull requests: %w", err)
	}
	return prs, nil
}

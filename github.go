package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh/v2"
	"time"
)

type GitHubService struct{}

type PullRequest struct {
	Author     Author     `json:"author"`
	CreatedAt  time.Time  `json:"createdAt"`
	IsDraft    bool       `json:"isDraft"`
	Number     int        `json:"number"`
	Repository Repository `json:"repository"`
	State      string     `json:"state"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
}

type Author struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type Repository struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
}

func (gs *GitHubService) GetOpenPullRequests() ([]PullRequest, error) {
	prs, _, err := gh.Exec("search", "prs", "--review-requested", "@me", "--state", "open", "--json", "repository,author,isDraft,createdAt,title,url,number")
	if err != nil {
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}
	var pullRequests []PullRequest

	json.Unmarshal([]byte(prs.String()), &pullRequests)
	return pullRequests, nil
}

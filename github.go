package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh/v2"
)

type GitHubService struct{}

type PullRequest struct {
	Author     Author     `json:"author"`
	IsDraft    bool       `json:"isDraft"`
	Repository Repository `json:"repository"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
}

type Author struct {
	Login string `json:"login"`
}

type Repository struct {
	NameWithOwner string `json:"nameWithOwner"`
}

func (gs *GitHubService) GetOpenPullRequests() ([]PullRequest, error) {
	prs, _, err := gh.Exec("search", "prs", "--review-requested", "@me", "--state", "open", "--json", "author,repository,isDraft,title,url")
	if err != nil {
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}
	var pullRequests []PullRequest

	json.Unmarshal([]byte(prs.String()), &pullRequests)
	return pullRequests, nil
}

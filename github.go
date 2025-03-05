package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh/v2"
	"strings"
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

func (gs *GitHubService) GetOpenPullRequests(cfg Config) ([]PullRequest, error) {
	prs, _, err := gh.Exec("search", "prs", "--review-requested", "@me", "--state", "open", "--json", "author,repository,isDraft,title,url")
	if err != nil {
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}
	var pullRequests []PullRequest

	if err = json.Unmarshal([]byte(prs.String()), &pullRequests); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pull requests: %w", err)
	}

	if cfg.Organization == BypassFilter {
		return pullRequests, nil
	}

	var filteredPRs []PullRequest
	organizations := strings.Split(cfg.Organization, "|")
	for _, pr := range pullRequests {
		for _, organization := range organizations {
			prOrg := strings.Split(pr.Repository.NameWithOwner, "/")[0]
			if prOrg == organization {
				filteredPRs = append(filteredPRs, pr)
				break
			}
		}
	}

	return filteredPRs, nil
}

package main

import (
	"bufio"
	"fmt"
	"github.com/cli/go-gh/v2"
	"strings"
)

type GitHubService struct{}

type PullRequest struct {
	Repo   string
	Id     string
	Status string
	Title  string
}

func (gs *GitHubService) GetOpenPullRequests() ([]PullRequest, error) {
	prs, _, err := gh.Exec("search", "prs", "--review-requested", "@me", "--state", "open")
	if err != nil {
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}
	return parsePullRequests(prs.String()), nil
}

func parsePullRequests(prLine string) []PullRequest {
	var prs []PullRequest
	scanner := bufio.NewScanner(strings.NewReader(prLine))
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		pr := PullRequest{
			Repo:   line[0],
			Id:     line[1],
			Status: line[2],
			Title:  line[3],
		}
		prs = append(prs, pr)
	}
	return prs
}

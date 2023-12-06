package main

import (
	"bufio"
	"fmt"
	"os/exec"
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
	cmd := exec.Command("gh", "search", "prs", "--review-requested", "@me", "--state", "open")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git status: %v", err)
	}
	return parsePullRequests(string(output)), nil
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

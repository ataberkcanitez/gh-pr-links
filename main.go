package main

import (
	"fmt"
	"os"
)

func main() {
	handleFlags()
	githubService := &GitHubService{}
	cliOutputHandler := &CliOutputHandler{}

	pullRequests, err := githubService.GetOpenPullRequests()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(pullRequests) == 0 {
		fmt.Println("No PR to show.")
		return
	}
	cliOutputHandler.PrintPullRequests(pullRequests)
}

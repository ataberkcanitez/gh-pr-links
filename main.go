package main

import (
	"fmt"
	"os"
)

func main() {
	options := handleOptions()
	githubService := &GitHubService{}

	pullRequests, err := githubService.GetOpenPullRequests(options)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(pullRequests) == 0 {
		fmt.Println("No PR to show.")
		return
	}

	cliOutputHandler := &CliOutputHandler{
		prs:    pullRequests,
		Config: options,
	}
	cliOutputHandler.PrintPullRequests()

}

// Command gh-pr-links is a GitHub CLI extension that lists the open pull
// requests awaiting your review, together with their links.
package main

import (
	"fmt"
	"os"

	"github.com/ataberkcanitez/gh-pr-links/internal/cli"
	"github.com/ataberkcanitez/gh-pr-links/internal/config"
	"github.com/ataberkcanitez/gh-pr-links/internal/github"
)

func main() {
	os.Exit(run())
}

func run() int {
	path, err := config.DefaultPath()
	if err != nil {
		fmt.Printf("Error: failed to read config file: %v\n", err)
		return cli.ExitError
	}

	app := &cli.App{
		Name:   os.Args[0],
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Config: config.NewStore(path),
		PRs:    github.NewClient(),
	}
	return app.Run(os.Args[1:])
}

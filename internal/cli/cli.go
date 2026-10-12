// Package cli implements the gh pr-links command.
package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ataberkcanitez/gh-pr-links/internal/config"
	"github.com/ataberkcanitez/gh-pr-links/internal/github"
	"github.com/ataberkcanitez/gh-pr-links/internal/render"
)

// Exit codes returned by Run.
const (
	ExitOK         = 0
	ExitError      = 1
	ExitUsageError = 2
)

// ConfigStore loads and saves the user's configuration.
type ConfigStore interface {
	Load() (config.Config, error)
	Save(config.Config) error
}

// PullRequestSource lists the pull requests awaiting the user's review.
type PullRequestSource interface {
	ReviewRequested() ([]github.PullRequest, error)
}

// App wires the command's dependencies together.
type App struct {
	// Name is the program name shown in flag usage errors.
	Name   string
	Stdout io.Writer
	// Stderr receives flag parsing errors and usage.
	Stderr io.Writer
	Config ConfigStore
	PRs    PullRequestSource
}

type options struct {
	help         bool
	style        string
	useEmoji     string
	organization string
}

func (a *App) parseFlags(args []string) (options, error) {
	var opts options
	fs := flag.NewFlagSet(a.Name, flag.ContinueOnError)
	fs.SetOutput(a.Stderr)

	styleUsage := "Sets the style of the output. Possible values: " + styleList()
	fs.BoolVar(&opts.help, "help", false, "Shows help message")
	fs.BoolVar(&opts.help, "h", false, "Shows help message")
	fs.StringVar(&opts.style, "style", "", styleUsage)
	fs.StringVar(&opts.style, "s", "", styleUsage)
	fs.StringVar(&opts.useEmoji, "use-emoji", "", "Use emoji in output")
	fs.StringVar(&opts.useEmoji, "u", "", "Use emoji in output")
	fs.StringVar(&opts.organization, "org", "", "Organization")
	fs.StringVar(&opts.organization, "o", "", "Organization")

	return opts, fs.Parse(args)
}

// Run executes the command with args (excluding the program name) and
// returns the process exit code.
//
// Setting --style or --use-emoji saves the preference and then lists pull
// requests; setting --org saves the filter and exits without listing.
func (a *App) Run(args []string) int {
	opts, err := a.parseFlags(args)
	if err != nil {
		// The flag package has already reported the error and usage.
		return ExitUsageError
	}

	cfg, err := a.Config.Load()
	if err != nil {
		fmt.Fprintf(a.Stdout, "Error: failed to read config file: %v\n", err)
		return ExitError
	}

	if opts.help {
		a.printHelp()
		return ExitOK
	}

	if opts.style != "" {
		if !render.IsValidStyle(opts.style) {
			fmt.Fprintf(a.Stdout, "Error: style value is invalid. Possible values: %s\n", styleList())
			return ExitError
		}
		cfg.Style = opts.style
		if err := a.Config.Save(cfg); err != nil {
			fmt.Fprintf(a.Stdout, "Error: failed to update config file: %v\n", err)
			return ExitError
		}
	}

	if opts.useEmoji != "" {
		if opts.useEmoji != "true" && opts.useEmoji != "false" {
			fmt.Fprintln(a.Stdout, "Error: use-emoji value is invalid. Possible values: true, false")
			return ExitError
		}
		cfg.UseEmoji = opts.useEmoji
		if err := a.Config.Save(cfg); err != nil {
			fmt.Fprintf(a.Stdout, "Error: failed to update config file: %v\n", err)
			return ExitError
		}
	}

	if opts.organization != "" {
		return a.setOrganization(cfg, opts.organization)
	}

	return a.listPullRequests(cfg)
}

func (a *App) setOrganization(cfg config.Config, org string) int {
	cfg.Organization = org
	if err := a.Config.Save(cfg); err != nil {
		fmt.Fprintf(a.Stdout, "Error: failed to update config: %v\n", err)
		return ExitError
	}

	if org == config.BypassFilter {
		fmt.Fprintln(a.Stdout, "Organization filter removed. All pull requests will be shown")
	} else {
		fmt.Fprintf(a.Stdout, "Organization filter set to '%s'. Only pull requests from these organizations will be shown.\n", org)
	}
	return ExitOK
}

func (a *App) listPullRequests(cfg config.Config) int {
	prs, err := a.PRs.ReviewRequested()
	if err != nil {
		fmt.Fprintln(a.Stdout, "Error:", err)
		return ExitError
	}

	if orgs := cfg.Organizations(); orgs != nil {
		prs = github.FilterByOrganizations(prs, orgs)
	}

	if len(prs) == 0 {
		fmt.Fprintln(a.Stdout, "No PR to show.")
		return ExitOK
	}

	table := render.Table{Style: cfg.Style, UseEmoji: cfg.EmojiEnabled()}
	if err := table.Render(a.Stdout, prs); err != nil {
		fmt.Fprintln(a.Stdout, "Error:", err)
		return ExitError
	}
	return ExitOK
}

func (a *App) printHelp() {
	w := a.Stdout
	fmt.Fprintln(w, "pr-links is a command line tool to show your open and review requested pull requests.")
	fmt.Fprintln(w, "Usage: gh pr-links [options]")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -h, --help\t\t\tShows help message")
	fmt.Fprintln(w, "  -s, --style <string>\t\tSets the style of the output. Possible values: "+styleList())
	fmt.Fprintln(w, "  -u, --use-emoji <bool>\tUse emoji in output. Possible values: true, false")
	fmt.Fprintln(w, "  -o, --org <string>\tFilter for organization, use '-' for skip the filtering, or <org_name> to filter based on organization")
	fmt.Fprintln(w, "\nConfiguration: gh pr-links --style=StyleCompactLite --use-emoji=true")
	fmt.Fprintln(w, "Usage: gh pr-links")
}

func styleList() string {
	return strings.Join(render.StyleNames(), ", ")
}

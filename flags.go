package main

import (
	"flag"
	"fmt"
	"os"
)

var styles = map[string]bool{
	"StyleCompactLite":    true,
	"StyleUnicode":        true,
	"StyleDefault":        true,
	"StyleCompact":        true,
	"StyleMarkdown":       true,
	"StyleRounded":        true,
	"StyleCompactClassic": true,
}

var (
	helpFlag         bool
	styleFlag        string
	useEmojiFlag     string
	organizationFlag string
)

func handleOptions() Config {
	flag.BoolVar(&helpFlag, "help", false, "Shows help message")
	flag.BoolVar(&helpFlag, "h", false, "Shows help message")

	flag.StringVar(&styleFlag, "style", "", "Sets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
	flag.StringVar(&styleFlag, "s", "", "Sets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")

	flag.StringVar(&useEmojiFlag, "use-emoji", "", "Use emoji in output")
	flag.StringVar(&useEmojiFlag, "u", "", "Use emoji in output")
	flag.StringVar(&organizationFlag, "org", "", "Organization")
	flag.StringVar(&organizationFlag, "o", "", "Organization")
	flag.Parse()

	config, err := ReadConfig()
	if err != nil {
		fmt.Printf("Error: failed to read config file: %v\n", err)
		os.Exit(1)
	}

	handleHelp(helpFlag)
	handleStyle(styleFlag, config)
	handleEmoji(useEmojiFlag, config)
	handleOrganization(organizationFlag, config)

	return *config
}

func handleOrganization(orgFlag string, config *Config) {
	if orgFlag == "" {
		return
	}
	config.Organization = orgFlag
	if err := UpdateConfig(config); err != nil {
		fmt.Printf("Error: failed to update config: %v\n", err)
		os.Exit(1)
	}
	if orgFlag == BypassFilter {
		fmt.Println("Organization filter removed. All pull requests will be shown")
	} else {
		fmt.Printf("Organization filter set to '%s'. Only pull requests from these organizations will be shown.\n", orgFlag)
	}
	os.Exit(0)
}

func handleHelp(help bool) {
	if help {
		printHelp()
		os.Exit(0)
	}
}

func handleStyle(selectedStyle string, config *Config) {
	if selectedStyle == "" {
		return
	}
	if _, ok := styles[selectedStyle]; !ok {
		fmt.Println("Error: style value is invalid. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
		os.Exit(1)
	}

	config.Style = selectedStyle
	if err := UpdateConfig(config); err != nil {
		fmt.Printf("Error: failed to update config file: %v\n", err)
		os.Exit(1)
	}
}

func handleEmoji(useEmoji string, config *Config) {
	if useEmoji == "" {
		return
	}

	if useEmoji != "true" && useEmoji != "false" {
		fmt.Println("Error: use-emoji value is invalid. Possible values: true, false")
		os.Exit(1)
	}

	config.UseEmoji = useEmoji
	if err := UpdateConfig(config); err != nil {
		fmt.Printf("Error: failed to update config file: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("pr-links is a command line tool to show your open and review requested pull requests.")
	fmt.Println("Usage: gh pr-links [options]")
	fmt.Println("Options:")
	fmt.Println("  -h, --help\t\t\tShows help message")
	fmt.Println("  -s, --style <string>\t\tSets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
	fmt.Println("  -u, --use-emoji <bool>\tUse emoji in output. Possible values: true, false")
	fmt.Println("  -o, --org <string>\tFilter for organization, use '-' for skip the filtering, or <org_name> to filter based on organization")
	fmt.Println("\nConfiguration: gh pr-links --style=StyleCompactLite --use-emoji=true")
	fmt.Println("Usage: gh pr-links")
}

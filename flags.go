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

func handleOptions() Config {
	helpFlag := flag.Bool("help", false, "Shows help message")
	styleFlag := flag.String("style", "", "Sets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
	useEmojiFlag := flag.String("use-emoji", "", "Use emoji in output")
	flag.Parse()

	config, err := ReadConfig()
	if err != nil {
		fmt.Printf("Error: failed to read config file: %v\n", err)
		os.Exit(1)
	}

	handleHelp(*helpFlag)
	handleStyle(*styleFlag, config)
	handleEmoji(*useEmojiFlag, config)

	return *config
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
	err := UpdateConfig(config)
	if err != nil {
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
	err := UpdateConfig(config)
	if err != nil {
		fmt.Printf("Error: failed to update config file: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("gh-pr is a command line tool to show your open and review requested pull requests.")
	fmt.Println("Usage: gh-pr [options]")
	fmt.Println("Options:")
	fmt.Println("  --help\t\tShows help message")
	fmt.Println("  --style <string>\tSets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
	fmt.Println("  --use-emoji <bool>\tUse emoji in output. Possible values: true, false")
	fmt.Println("\nConfiguration: gh pr-links --style=StyleCompactLite --use-emoji=true")
	fmt.Println("Usage: gh pr-links")
}

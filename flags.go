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

func handleFlags() {
	helpFlag := flag.Bool("help", false, "Shows help message")
	styleFlag := flag.String("style", "", "Sets the style of the output. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
	useEmojiFlag := flag.String("use-emoji", "", "Use emoji in output")
	flag.Parse()

	handleHelp(*helpFlag)
	handleStyle(*styleFlag)
	handleEmoji(*useEmojiFlag)

}

func handleHelp(help bool) {
	if help {
		printHelp()
		os.Exit(0)
	}

}

func handleStyle(selectedStyle string) {
	if selectedStyle == "" {
		return
	}
	if _, ok := styles[selectedStyle]; !ok {
		fmt.Println("Error: style value is invalid. Possible values: StyleCompactLite, StyleUnicode, StyleDefault, StyleCompact, StyleMarkdown, StyleRounded, StyleCompactClassic")
		os.Exit(1)
	}

	err := os.Setenv("GH_PR_STYLE", selectedStyle)
	if err != nil {
		fmt.Printf("Error: failed to set GH_PR_STYLE environment variable: %v\n", err)
		return
	}
}

func handleEmoji(useEmoji string) {
	if useEmoji == "" {
		return
	}
	switch useEmoji {
	case "true":
		err := os.Setenv("GH_PR_USE_EMOJI", "true")
		if err != nil {
			fmt.Printf("Error: failed to set GH_PR_USE_EMOJI environment variable: %v\n", err)
			return
		}
	case "false":
		err := os.Setenv("GH_PR_USE_EMOJI", "false")
		if err != nil {
			fmt.Printf("Error: failed to set GH_PR_USE_EMOJI environment variable: %v\n", err)
			return
		}
	default:
		fmt.Println("Error: use-emoji value is invalid. Possible values: true, false")
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

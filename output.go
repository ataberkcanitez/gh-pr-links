package main

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"strings"
)

type CliOutputHandler struct {
	prs    []PullRequest
	Config Config
}

func (co *CliOutputHandler) PrintPullRequests() {
	table := simpletable.New()
	table.Header = prepareHeader()
	for _, pr := range co.prs {
		title := co.prepareTitle(pr.Title, pr.IsDraft)
		author := prepareAuthor(pr.Author.Login, pr.IsDraft)
		url := prepareUrl(pr.URL, pr.IsDraft)
		owner := prepareOwner(pr.Repository.NameWithOwner, pr.IsDraft)
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: owner},
			{Align: simpletable.AlignLeft, Text: title},
			{Align: simpletable.AlignLeft, Text: author},
			{Align: simpletable.AlignLeft, Text: url},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	style := co.getStyle()
	table.SetStyle(style)
	fmt.Println(table.String())
}

func (co *CliOutputHandler) getStyle() *simpletable.Style {
	selectedStyle := co.Config.Style
	switch selectedStyle {
	case "StyleCompactLite":
		return simpletable.StyleCompactLite
	case "StyleUnicode":
		return simpletable.StyleUnicode
	case "StyleDefault":
		return simpletable.StyleDefault
	case "StyleCompact":
		return simpletable.StyleCompact
	case "StyleMarkdown":
		return simpletable.StyleMarkdown
	case "StyleRounded":
		return simpletable.StyleRounded
	case "StyleCompactClassic":
		return simpletable.StyleCompactClassic
	default:
		return simpletable.StyleRounded
	}
}

func prepareOwner(owner string, draft bool) string {
	textColor := green
	if draft {
		textColor = gray
	}
	return textColor(owner)
}

func prepareUrl(url string, draft bool) string {
	urlColor := blue
	if draft {
		urlColor = gray
	}
	return urlColor(url)
}

func prepareAuthor(login string, draft bool) string {
	textColor := green
	if draft {
		textColor = gray
	}
	return textColor(login)
}

func prepareHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "REPO"},
			{Align: simpletable.AlignCenter, Text: "TITLE"},
			{Align: simpletable.AlignCenter, Text: "AUTHOR"},
			{Align: simpletable.AlignCenter, Text: "LINK"},
		},
	}
}

func (co *CliOutputHandler) prepareTitle(input string, isDraft bool) string {
	textColor := green
	if isDraft {
		textColor = gray
	}
	emoji := co.getEmoji(isDraft)

	const maxCharacters = 40
	var result strings.Builder
	var count int
	words := strings.Fields(input)

	useEmoji := co.Config.UseEmoji
	if useEmoji == "true" {
		result.WriteString(emoji)
		result.WriteString(" ")
		count += 2
	}

	for _, word := range words {
		wordLen := len(word)
		if count+wordLen > maxCharacters {
			result.WriteString("\n")
			count = 0
		}
		result.WriteString(textColor(word))
		result.WriteString(" ")
		count += wordLen + 1
	}

	return result.String()
}
func (co *CliOutputHandler) getEmoji(isDraft bool) string {
	useEmoji := co.Config.UseEmoji
	if useEmoji == "false" {
		return ""
	}

	if isDraft {
		return "\U0001F527"
	}
	return "\U0001F44C"
}

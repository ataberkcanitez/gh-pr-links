package main

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"strings"
)

type CliOutputHandler struct{}

func (co *CliOutputHandler) PrintPullRequests(prs []PullRequest) {
	table := simpletable.New()
	table.Header = prepareHeader()
	for _, row := range prs {
		link := prepareLink(row.Repo, row.Id)
		title := prepareTitle(row.Title)
		r := []*simpletable.Cell{
			{Text: row.Repo},
			{Text: row.Id},
			{Text: title},
			{Text: link},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleUnicode)
	fmt.Println(table.String())
}

func prepareHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "REPO"},
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "TITLE"},
			{Align: simpletable.AlignCenter, Text: "LINK"},
		},
	}
}

func prepareLink(repo, id string) string {
	return fmt.Sprintf("https://www.github.com/%s/pull/%s", repo, id)
}

func prepareTitle(input string) string {
	const maxCharacters = 40
	var result strings.Builder
	var count int
	words := strings.Fields(input)

	for _, word := range words {
		wordLen := len(word)
		if count+wordLen > maxCharacters {
			result.WriteString("\n")
			count = 0
		}
		result.WriteString(word)
		result.WriteString(" ")
		count += wordLen + 1
	}

	return result.String()
}

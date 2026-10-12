// Package render draws pull requests as a colored terminal table.
package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/alexeyco/simpletable"

	"github.com/ataberkcanitez/gh-pr-links/internal/github"
)

const (
	// maxTitleWidth is the approximate number of characters after which a
	// title wraps onto a new line.
	maxTitleWidth = 40

	readyEmoji = "\U0001F44C" // 👌
	draftEmoji = "\U0001F527" // 🔧
)

// Table renders pull requests as a table.
type Table struct {
	// Style is a style name as accepted by IsValidStyle. Unknown or empty
	// names fall back to DefaultStyle.
	Style string
	// UseEmoji prefixes each title with an emoji indicating draft status.
	UseEmoji bool
}

// Render writes prs to w as a table followed by a newline.
func (t Table) Render(w io.Writer, prs []github.PullRequest) error {
	table := simpletable.New()
	table.Header = header()
	for _, pr := range prs {
		table.Body.Cells = append(table.Body.Cells, t.row(pr))
	}
	table.SetStyle(lookupStyle(t.Style))

	_, err := fmt.Fprintln(w, table.String())
	return err
}

func header() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "REPO"},
			{Align: simpletable.AlignCenter, Text: "TITLE"},
			{Align: simpletable.AlignCenter, Text: "AUTHOR"},
			{Align: simpletable.AlignCenter, Text: "LINK"},
		},
	}
}

func (t Table) row(pr github.PullRequest) []*simpletable.Cell {
	text, link := green, blue
	if pr.IsDraft {
		text, link = gray, gray
	}

	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: text(pr.Repository.NameWithOwner)},
		{Align: simpletable.AlignLeft, Text: t.title(pr.Title, pr.IsDraft, text)},
		{Align: simpletable.AlignLeft, Text: text(pr.Author.Login)},
		{Align: simpletable.AlignLeft, Text: link(pr.URL)},
	}
}

// title colors each word of s and wraps it at roughly maxTitleWidth
// characters, optionally prefixing it with a draft-status emoji.
func (t Table) title(s string, isDraft bool, color func(string) string) string {
	var b strings.Builder
	width := 0

	if t.UseEmoji {
		emoji := readyEmoji
		if isDraft {
			emoji = draftEmoji
		}
		b.WriteString(emoji)
		b.WriteString(" ")
		width += 2
	}

	for word := range strings.FieldsSeq(s) {
		// Width is measured in bytes, matching the original layout.
		if width+len(word) > maxTitleWidth {
			b.WriteString("\n")
			width = 0
		}
		b.WriteString(color(word))
		b.WriteString(" ")
		width += len(word) + 1
	}
	return b.String()
}

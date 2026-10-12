package render

import "github.com/alexeyco/simpletable"

// DefaultStyle is used when no style or an unknown style is configured.
const DefaultStyle = "StyleRounded"

// styleNames lists the supported styles in the order they are shown to users.
var styleNames = []string{
	"StyleCompactLite",
	"StyleUnicode",
	"StyleDefault",
	"StyleCompact",
	"StyleMarkdown",
	"StyleRounded",
	"StyleCompactClassic",
}

var styles = map[string]*simpletable.Style{
	"StyleCompactLite":    simpletable.StyleCompactLite,
	"StyleUnicode":        simpletable.StyleUnicode,
	"StyleDefault":        simpletable.StyleDefault,
	"StyleCompact":        simpletable.StyleCompact,
	"StyleMarkdown":       simpletable.StyleMarkdown,
	"StyleRounded":        simpletable.StyleRounded,
	"StyleCompactClassic": simpletable.StyleCompactClassic,
}

// StyleNames returns the names of all supported table styles.
func StyleNames() []string {
	return append([]string(nil), styleNames...)
}

// IsValidStyle reports whether name is a supported table style.
func IsValidStyle(name string) bool {
	_, ok := styles[name]
	return ok
}

func lookupStyle(name string) *simpletable.Style {
	if s, ok := styles[name]; ok {
		return s
	}
	return styles[DefaultStyle]
}

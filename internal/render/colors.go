package render

// ANSI escape sequences used to color table cells.
const (
	colorDefault = "\x1b[39m"
	colorGreen   = "\x1b[32m"
	colorBlue    = "\x1b[94m"
	colorGray    = "\x1b[90m"
)

func green(s string) string { return colorGreen + s + colorDefault }
func blue(s string) string  { return colorBlue + s + colorDefault }
func gray(s string) string  { return colorGray + s + colorDefault }

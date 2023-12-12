package main

import "fmt"

const (
	ColorDefault = "\x1b[39m"
	ColorGreen   = "\x1b[32m"
	ColorBlue    = "\x1b[94m"
	ColorGray    = "\x1b[90m"
)

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func gray(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGray, s, ColorDefault)
}

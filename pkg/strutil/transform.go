package strutil

import (
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// RightPad pads str with spaces on the right up to spaceCount runes.
// If the trimmed string is already as long as or longer than spaceCount,
// a single trailing space is added.
func RightPad(str string, spaceCount uint8) string {
	str = strings.TrimSpace(str)
	padCount := int(spaceCount) - utf8.RuneCountInString(str)
	if padCount < 1 {
		return str + " "
	}
	return str + strings.Repeat(" ", padCount)
}

// UnicodeTruncate performs unicode-aware string truncation.
// When truncation occurs, it appends "…".
func UnicodeTruncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if uniseg.StringWidth(s) <= width {
		return s
	}
	var result strings.Builder
	remaining := width - uniseg.StringWidth("…")
	graphemes := uniseg.NewGraphemes(s)
	for graphemes.Next() {
		if graphemes.Width() > remaining {
			break
		}
		result.WriteString(graphemes.Str())
		remaining -= graphemes.Width()
	}
	result.WriteRune('…')
	return result.String()
}

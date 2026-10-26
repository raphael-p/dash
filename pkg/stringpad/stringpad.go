package stringpad

import (
	"strings"
	"unicode/utf8"
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

// Rune-wise string truncation, when truncation occurs, appends …
func TruncateString(s string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

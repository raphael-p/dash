package tviewcomponents

import "github.com/gdamore/tcell/v2"

var defaultHighlightColour tcell.Color
var defaultFallback KeybindFallback = func(_ string) {}

func SetHighlightColour(highlightColour tcell.Color) {
	defaultHighlightColour = highlightColour
}

func SetKeybindFallback(fallback KeybindFallback) {
	defaultFallback = fallback
}

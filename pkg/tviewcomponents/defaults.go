package tviewcomponents

import "github.com/gdamore/tcell/v2"

var defaultHighlightColour, defaultSecondaryColour tcell.Color
var defaultFallback KeybindFallback = func(_ string) {}

func SetColours(highlightColour, secondaryColour tcell.Color) {
	defaultHighlightColour = highlightColour
	defaultSecondaryColour = secondaryColour
}

func SetKeybindFallback(fallback KeybindFallback) {
	defaultFallback = fallback
}

package keybindmenu

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Bind = func(eventKey tcell.Key, key rune) bool

var (
	DefaultBind   Bind = func(_ tcell.Key, _ rune) bool { return false } // just binds on Keybind.key
	BindEnter     Bind = func(e tcell.Key, _ rune) bool { return e == tcell.KeyEnter }
	BindBackspace Bind = func(e tcell.Key, _ rune) bool { return e == tcell.KeyBackspace || e == tcell.KeyBackspace2 }
	BindNumber    Bind = func(e tcell.Key, r rune) bool { return e == tcell.KeyRune && r >= '0' && r <= '9' }
)

type Fallback = func(commandName string)
type KeybindMenu struct {
	fallback        Fallback
	footnote        string
	highlightColour tcell.Color
	keybinds        []Keybind
	isScrollable    bool
}

type Handler = func(key rune)
type DisableKeybind = func() bool
type Keybind struct {
	handler     Handler
	bind        Bind
	description string
	key         rune
}

func New() *KeybindMenu {
	return &KeybindMenu{func(_ string) {}, "", 0, []Keybind{}, true}
}

func (hm *KeybindMenu) SetFallback(fallback Fallback) *KeybindMenu {
	hm.fallback = fallback
	return hm
}

func (hm *KeybindMenu) SetFootnote(footnote string) *KeybindMenu {
	hm.footnote = footnote
	return hm
}

func (hm *KeybindMenu) SetHighlighColour(colour tcell.Color) *KeybindMenu {
	hm.highlightColour = colour
	return hm
}

func (hm *KeybindMenu) AddKeybind(key rune, description string, bind Bind, handler Handler) *KeybindMenu {
	hm.keybinds = append(hm.keybinds, Keybind{handler, bind, description, key})
	return hm
}

func (hm *KeybindMenu) DisableScrolling() *KeybindMenu {
	hm.isScrollable = false
	return hm
}

type applyable interface {
	SetText(string) *tview.TextView
	SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
}

func (hm *KeybindMenu) Apply(parentComponent applyable) {
	parentComponent.SetText(hm.generateText())
	parentComponent.SetInputCapture(hm.generateInputCapture())
}

func (hm *KeybindMenu) generateText() string {
	builder := new(strings.Builder)

	for _, keybind := range hm.keybinds {
		hasKey := keybind.key != 0
		hasDescription := keybind.description != ""
		if builder.Len() != 0 && (hasKey || hasDescription) {
			fmt.Fprintf(builder, "\n")
		}

		if hasKey {
			fmt.Fprintf(builder, "([%s::b]%c[-:-:-])", hm.highlightColour, keybind.key)
		}

		if hasDescription {
			fmt.Fprintf(builder, "\t%s", keybind.description)
		}
	}

	if hm.footnote != "" {
		if builder.Len() != 0 {
			fmt.Fprintf(builder, "\n")
		}
		fmt.Fprintf(builder, "[::d]%s[::-]", hm.footnote)
	}

	return builder.String()
}

func (hm *KeybindMenu) generateInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		eventKey := event.Key()
		key := event.Rune()

		if hm.isScrollable && (eventKey == tcell.KeyDown || eventKey == tcell.KeyUp) {
			return event
		}

		hasTriggered := false
		for _, keybind := range hm.keybinds {
			if (key != 0 && key == keybind.key) || keybind.bind(eventKey, key) {
				keybind.handler(key)
				hasTriggered = true
				break
			}
		}

		if !hasTriggered {
			if key == 0 || (string(key) != " " && unicode.IsSpace(key)) {
				hm.fallback(event.Name())
			} else {
				hm.fallback("'" + string(key) + "'")
			}
		}

		return nil
	}
}

package keybindmenu

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Bind = func(eventKey tcell.Key, key rune, mod tcell.ModMask) bool

var (
	DefaultBind   Bind = func(_ tcell.Key, _ rune, _ tcell.ModMask) bool { return false } // just binds on Keybind.key
	BindBackspace Bind = func(e tcell.Key, _ rune, _ tcell.ModMask) bool {
		return e == tcell.KeyBackspace || e == tcell.KeyBackspace2
	}
	BindClearLine Bind = func(_ tcell.Key, r rune, m tcell.ModMask) bool {
		return ((r == 'u' || r == 'U') && m == tcell.ModCtrl)
	}
	BindEnter Bind = func(e tcell.Key, _ rune, m tcell.ModMask) bool {
		return e == tcell.KeyEnter && m == tcell.ModNone
	}
	BindNewLine Bind = func(e tcell.Key, r rune, m tcell.ModMask) bool {
		return (e == tcell.KeyEnter && m == tcell.ModShift) || ((r == 'j' || r == 'J') && m == tcell.ModCtrl)
	}
	BindNumber Bind = func(e tcell.Key, r rune, _ tcell.ModMask) bool {
		return e == tcell.KeyRune && r >= '0' && r <= '9'
	}
	BindPrintable Bind = func(e tcell.Key, r rune, m tcell.ModMask) bool {
		return e == tcell.KeyRune
	}
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

var defaultHighlightColour tcell.Color
var defaultFallback Fallback = func(_ string) {}

func New() *KeybindMenu {
	return &KeybindMenu{defaultFallback, "", defaultHighlightColour, []Keybind{}, true}
}

func SetDefaults(fallback Fallback, highlightColour tcell.Color) {
	defaultFallback = fallback
	defaultHighlightColour = highlightColour
}

func (hm *KeybindMenu) SetFootnote(footnote string) *KeybindMenu {
	hm.footnote = footnote
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

func (hm *KeybindMenu) Text() string {
	builder := new(strings.Builder)

	descriptionOnly := []Keybind{}
	for _, keybind := range hm.keybinds {
		hasKey := keybind.key != 0
		hasDescription := keybind.description != ""
		if hasKey {
			if builder.Len() != 0 {
				fmt.Fprintf(builder, "\n")
			}
			fmt.Fprintf(builder, "([%s::b]%c[-:-:-])", hm.highlightColour, keybind.key)
			if hasDescription {
				fmt.Fprintf(builder, "\t%s", keybind.description)
			}
		} else if hasDescription {
			descriptionOnly = append(descriptionOnly, keybind)
		}
	}

	// display keybinds with a description but no keys in a separate block
	for i, keybind := range descriptionOnly {
		if builder.Len() != 0 {
			fmt.Fprintf(builder, "\n")
			if i == 0 {
				fmt.Fprintf(builder, "\n")
			}
		}
		fmt.Fprintf(builder, "%s", keybind.description)
	}

	if hm.footnote != "" {
		if builder.Len() != 0 {
			fmt.Fprintf(builder, "\n")
		}
		fmt.Fprintf(builder, "[::d]%s[::-]", hm.footnote)
	}

	return builder.String()
}

func (hm *KeybindMenu) InputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		eventKey := event.Key()
		key := event.Rune()
		mods := event.Modifiers()

		if hm.isScrollable && (eventKey == tcell.KeyDown || eventKey == tcell.KeyUp || eventKey == tcell.KeyLeft || eventKey == tcell.KeyRight) {
			return event
		}

		hasTriggered := false
		for _, keybind := range hm.keybinds {
			if (key != 0 && key == keybind.key && mods == tcell.ModNone) || keybind.bind(eventKey, key, mods) {
				// intentionally allow default behaviour
				if keybind.handler == nil {
					return event
				}
				keybind.handler(key)
				hasTriggered = true
				break
			}
		}

		if !hasTriggered {
			eventName := event.Name()
			if strings.HasPrefix(eventName, "Rune") {
				hm.fallback("'" + string(key) + "'")
			} else {
				hm.fallback(event.Name())
			}
		}

		return nil
	}
}

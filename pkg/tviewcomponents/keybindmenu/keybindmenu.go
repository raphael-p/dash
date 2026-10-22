package keybindmenu

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Bind = func(eventKey tcell.Key, key rune, mod tcell.ModMask) bool

var (
	DefaultBind Bind = func(_ tcell.Key, _ rune, _ tcell.ModMask) bool { return false } // just binds on Keybind.key
	BindEnter   Bind = func(e tcell.Key, _ rune, m tcell.ModMask) bool {
		return e == tcell.KeyEnter && m == tcell.ModNone
	}
	BindBackspace Bind = func(e tcell.Key, _ rune, _ tcell.ModMask) bool {
		return e == tcell.KeyBackspace || e == tcell.KeyBackspace2
	}
	BindNewLine Bind = func(e tcell.Key, _ rune, m tcell.ModMask) bool {
		return (e == tcell.KeyEnter && m == tcell.ModShift)
	}
	BindNumber Bind = func(e tcell.Key, r rune, _ tcell.ModMask) bool {
		return e == tcell.KeyRune && r >= '0' && r <= '9'
	}
	BindPrintable Bind = func(e tcell.Key, r rune, _ tcell.ModMask) bool {
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

func (hm *KeybindMenu) SetHighlightColour(colour tcell.Color) *KeybindMenu {
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

func (hm *KeybindMenu) GenerateText() string {
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

func (hm *KeybindMenu) GenerateInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		// map CTRL+J to SHIFT+ENTER to support newlines
		if (event.Rune() == 'j' || event.Rune() == 'J') && event.Modifiers() == tcell.ModCtrl {
			event = tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModShift)
		}

		eventKey := event.Key()
		key := event.Rune()
		mods := event.Modifiers()

		if hm.isScrollable && (eventKey == tcell.KeyDown || eventKey == tcell.KeyUp) {
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

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
	BindNumber    Bind = func(_ tcell.Key, r rune) bool { return r >= '0' && r <= '9' }
)

type DefaultKeybindAction = func(commandName string)

type KeybindAction = func(key rune)

type KeybindMenu struct {
	defaultAction   DefaultKeybindAction
	footnote        string
	highlightColour tcell.Color
	keybinds        []Keybind
	isScrollable    bool
}

type Keybind struct {
	key         rune
	description string
	bind        Bind
	callback    func(key rune)
}

func New() *KeybindMenu {
	return &KeybindMenu{func(_ string) {}, "", 0, []Keybind{}, true}
}

func (hm *KeybindMenu) SetDefaultAction(defaultAction DefaultKeybindAction) *KeybindMenu {
	hm.defaultAction = defaultAction
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

func (hm *KeybindMenu) AddKeybind(key rune, description string, trigger Bind, callback KeybindAction) *KeybindMenu {
	hm.keybinds = append(hm.keybinds, Keybind{key, description, trigger, callback})
	return hm
}

func (hm *KeybindMenu) DisableScrolling() *KeybindMenu {
	hm.isScrollable = false
	return hm
}

func (hm *KeybindMenu) Apply(parentComponent *tview.TextView) {
	parentComponent.SetText(hm.generateText())
	parentComponent.SetInputCapture(hm.generateInputCapture())
}

func (hm *KeybindMenu) generateText() string {
	var builder strings.Builder
	for idx, keybind := range hm.keybinds {
		if keybind.key == 0 {
			continue
		}

		builder.WriteString(fmt.Sprintf(
			"([%s::b]%c[-:-:-]) %s",
			hm.highlightColour, keybind.key, keybind.description,
		))

		if hm.footnote != "" || idx < len(hm.keybinds)-1 {
			builder.WriteString("\n")
		}
	}

	builder.WriteString(fmt.Sprintf("[::d]%s[::-]", hm.footnote))

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
			if key == keybind.key || keybind.bind(eventKey, key) {
				keybind.callback(key)
				hasTriggered = true
				break
			}
		}

		if !hasTriggered {
			if key == 0 || (string(key) != " " && unicode.IsSpace(key)) {
				hm.defaultAction(event.Name())
			} else {
				hm.defaultAction("'" + string(key) + "'")
			}
		}

		return nil
	}
}

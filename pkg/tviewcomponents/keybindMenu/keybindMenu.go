package keybindmenu

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type defaultAction = func(eventName string, eventRune rune)

type KeybindMenu struct {
	defaultAction   defaultAction
	footnote        string
	highlightColour tcell.Color
	keybinds        []Keybind
	isScrollable    bool
}

type trigger = func(eventKey tcell.Key, key rune) bool
type callback = func(key rune)

type Keybind struct {
	key         rune
	description string
	trigger     trigger
	callback    func(key rune)
}

func New() *KeybindMenu {
	return &KeybindMenu{func(e string, r rune) {}, "", 0, []Keybind{}, true}
}

func (hm *KeybindMenu) SetDefaultAction(defaultAction defaultAction) *KeybindMenu {
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

func (hm *KeybindMenu) AddKeybind(key rune, description string, trigger trigger, callback callback) *KeybindMenu {
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
			if key == keybind.key || keybind.trigger(eventKey, key) {
				keybind.callback(key)
				hasTriggered = true
				break
			}
		}

		if !hasTriggered {
			hm.defaultAction(event.Name(), key)
		}

		return nil
	}
}

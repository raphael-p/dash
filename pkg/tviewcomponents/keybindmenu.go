package tviewcomponents

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/pkg/strutil"
)

type Bind = func(eventKey tcell.Key, key rune, mod tcell.ModMask) bool

var (
	BindBackspace Bind = func(e tcell.Key, _ rune, _ tcell.ModMask) bool {
		return e == tcell.KeyBackspace || e == tcell.KeyBackspace2
	}
	BindCtrlJ Bind = func(e tcell.Key, r rune, m tcell.ModMask) bool {
		return (e == tcell.KeyEnter && m == tcell.ModShift) || ((r == 'j' || r == 'J') && m == tcell.ModCtrl)
	}
	BindCtrlU Bind = func(_ tcell.Key, r rune, m tcell.ModMask) bool {
		return ((r == 'u' || r == 'U') && m == tcell.ModCtrl)
	}
	BindEnter Bind = func(e tcell.Key, _ rune, m tcell.ModMask) bool {
		return e == tcell.KeyEnter && m == tcell.ModNone
	}
	BindEsc Bind = func(e tcell.Key, _ rune, _ tcell.ModMask) bool {
		return e == tcell.KeyEscape
	}
	BindNumber Bind = func(e tcell.Key, r rune, _ tcell.ModMask) bool {
		return e == tcell.KeyRune && r >= '0' && r <= '9'
	}
	BindPrintable Bind = func(e tcell.Key, r rune, m tcell.ModMask) bool {
		return e == tcell.KeyRune
	}
	BindRune Bind = func(_ tcell.Key, _ rune, _ tcell.ModMask) bool { return false } // just binds on Keybind.key
	BindTab  Bind = func(e tcell.Key, _ rune, _ tcell.ModMask) bool {
		return e == tcell.KeyTab || e == tcell.KeyTAB
	}
)

type orderedKeybinds struct {
	runeKeybinds, primaryKeybinds, secondaryKeybinds, passthroughKeybinds []Keybind
}

func (kg orderedKeybinds) getAll() []Keybind {
	return slices.Concat(kg.runeKeybinds, kg.primaryKeybinds, kg.secondaryKeybinds, kg.passthroughKeybinds)
}

func (kg orderedKeybinds) getVisible() []Keybind {
	return slices.Concat(kg.runeKeybinds, kg.primaryKeybinds, kg.secondaryKeybinds)
}

type KeybindFallback = func(commandName string)
type KeybindMenu struct {
	fallback     KeybindFallback
	keybinds     orderedKeybinds
	isScrollable bool
}

type KeybindHandler = func()
type DisableKeybind = func() bool
type Keybind struct {
	handler KeybindHandler
	bind    Bind
	display string
	key     rune
}

func NewKeybindMenu() *KeybindMenu {
	return &KeybindMenu{defaultFallback, orderedKeybinds{}, true}
}

func (hm *KeybindMenu) AddRuneKeybind(key rune, description string, handler KeybindHandler) *KeybindMenu {
	display := fmt.Sprintf("([%s::b]%c[-:-:-])  %s", defaultHighlightColour, key, description)
	hm.keybinds.runeKeybinds = append(hm.keybinds.runeKeybinds, Keybind{handler, BindRune, display, key})
	return hm
}

func (hm *KeybindMenu) AddPrimaryKeybind(shortcut, function string, bind Bind, handler KeybindHandler) *KeybindMenu {
	display := fmt.Sprintf("[%s::b]%s[-:-:-]%s", defaultHighlightColour, strutil.RightPad(shortcut, 8), function)
	hm.keybinds.primaryKeybinds = append(hm.keybinds.primaryKeybinds, Keybind{handler, bind, display, 0})
	return hm
}

func (hm *KeybindMenu) AddSecondaryKeybind(shortcut, function string, bind Bind, handler KeybindHandler) *KeybindMenu {
	display := fmt.Sprintf("[%s::b]%s[-:-:-]%s", defaultSecondaryColour, strutil.RightPad(shortcut, 8), function)
	hm.keybinds.secondaryKeybinds = append(hm.keybinds.secondaryKeybinds, Keybind{handler, bind, display, 0})
	return hm
}

func (hm *KeybindMenu) AddPassthroughKeybind(bind Bind) *KeybindMenu {
	hm.keybinds.passthroughKeybinds = append(hm.keybinds.passthroughKeybinds, Keybind{nil, bind, "", 0})
	return hm
}

func (hm *KeybindMenu) DisableScrolling() *KeybindMenu {
	hm.isScrollable = false
	return hm
}

func (hm *KeybindMenu) Text() string {
	builder := new(strings.Builder)

	numRuneKeybinds := len(hm.keybinds.runeKeybinds)
	for idx, keybind := range hm.keybinds.getVisible() {
		if keybind.display == "" {
			continue
		}
		if builder.Len() != 0 {
			builder.WriteRune('\n')
			if idx == numRuneKeybinds {
				builder.WriteRune('\n')
			}
		}
		builder.WriteString(keybind.display)
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
		for _, keybind := range hm.keybinds.getAll() {
			if (key != 0 && key == keybind.key && mods == tcell.ModNone) || keybind.bind(eventKey, key, mods) {
				// intentionally allow default behaviour
				if keybind.handler == nil {
					return event
				}
				keybind.handler()
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

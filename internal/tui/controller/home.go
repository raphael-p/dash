package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.refreshTasks()

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)

	handler := &keybindHandler{c}
	hasInput := func() bool { return c.infoPanel.GetInput() != "" }
	fallback := func(commandName string) { c.infoPanel.Warn("invalid command: " + commandName) }

	setWrapper := func() { setHomeKeybinds(home, handler, hasInput, fallback) }
	c.infoPanel.SetOnInputChange(func(_ string) { setWrapper() })
	setWrapper()

	c.inputPanel.Set("Welcome to Dash!", home)
}

func setHomeKeybinds(home *tview.TextView, handler *keybindHandler, hasInput func() bool, fallback func(string)) {
	menu := keybindmenu.New().SetHighlighColour(tcell.ColorLimeGreen).SetFallback(fallback)

	if hasInput() {
		menu.AddKeybind('o', "open task", keybindmenu.BindEnter, handler.open)
		menu.AddKeybind('b', "bump task priority", keybindmenu.DefaultBind, handler.bump)
		menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.remove)
		menu.AddKeybind(0, "", keybindmenu.BindNumber, handler.numberInput)
		menu.AddKeybind(0, "press backspace to clear", keybindmenu.BindBackspace, handler.backspace)
	} else {
		menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.dash)
		menu.AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.add)
		menu.AddKeybind(0, "enter a task ID", keybindmenu.BindNumber, handler.numberInput)
		menu.AddKeybind(0, "", keybindmenu.BindBackspace, handler.backspace)
	}

	menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDown).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUp).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit).
		Apply(home)
}

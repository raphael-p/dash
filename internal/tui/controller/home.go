package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.refreshTasks()

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)

	handler := &keybindHandler{c}
	hasInput := func() bool { return c.infoPanel.GetInput() != "" }

	setWrapper := func() { setHomeKeybinds(home, handler, hasInput, c.displayPanel.ShowCompleted) }
	c.infoPanel.SetOnInputChange(func(_ string) { setWrapper() })
	setWrapper()

	c.inputPanel.Set("Welcome to Dash!", home)
}

func setHomeKeybinds(home *tview.TextView, handler *keybindHandler, hasInput func() bool, showCompleted func() bool) {
	menu := keybindmenu.New().
		SetHighlighColour(tcell.ColorLimeGreen).
		SetFallback(handler.fallback)

	if hasInput() {
		menu.AddKeybind('o', "open task", keybindmenu.BindEnter, handler.openTask).
			AddKeybind('c', "copy task to clipboard", keybindmenu.DefaultBind, handler.copyTask)
		if !showCompleted() {
			menu.AddKeybind('b', "bump task priority", keybindmenu.DefaultBind, handler.bumpTask)
		}
		menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTask).
			AddKeybind(0, "", keybindmenu.BindNumber, handler.numberInput).
			AddKeybind(0, "press backspace to clear", keybindmenu.BindBackspace, handler.backspace)
	} else {
		if showCompleted() {
			menu.AddKeybind('t', "show active tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		} else {
			menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.startDash).
				AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addFromHome).
				AddKeybind('t', "show completed tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		}
		menu.AddKeybind(0, "enter a task ID", keybindmenu.BindNumber, handler.numberInput).
			AddKeybind(0, "", keybindmenu.BindBackspace, handler.backspace)
	}

	menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDown).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUp).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit).
		Apply(home)
}

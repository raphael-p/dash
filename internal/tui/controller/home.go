package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/pkg/tviewcomponents/commandline"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.refreshTasks()

	commandLine := commandline.New().SetChevronColour(tcell.ColorLimeGreen).SetPlaceholder("type a number...")
	commandLine.SetInput("")
	keybindMenuContainer := tview.NewTextView().SetDynamicColors(true)

	home := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandLine.GetTextView(), 0, 1, false).
		AddItem(keybindMenuContainer, 0, 5, true)
	home.SetBorderPadding(1, 1, 2, 2)

	handler := &keybindHandler{c}
	hasInput := func() bool { return commandLine.GetInput() != "" }

	setWrapper := func() { setHomeKeybinds(keybindMenuContainer, handler, hasInput, c.displayPanel.ShowCompleted) }
	commandLine.SetOnInputChange(func(_ string) { setWrapper() })
	setWrapper()

	c.inputPanel.Set("Welcome to Dash!", home, commandLine.GetInput, commandLine.SetInput)
}

func setHomeKeybinds(home *tview.TextView, handler *keybindHandler, hasInput func() bool, showCompleted func() bool) {
	menu := keybindmenu.New().
		SetHighlightColour(tcell.ColorLimeGreen).
		SetFallback(handler.fallback)

	if hasInput() {
		menu.AddKeybind('o', "open task", keybindmenu.BindEnter, handler.openTask).
			AddKeybind('c', "copy task to clipboard", keybindmenu.DefaultBind, handler.copyTask)
		if showCompleted() {
			menu.AddKeybind('u', "reactivate task", keybindmenu.DefaultBind, handler.reactivateTask)
		} else {
			menu.AddKeybind('b', "bump task priority", keybindmenu.DefaultBind, handler.bumpTask)
		}
		menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTask).
			AddKeybind(0, "", keybindmenu.BindNumber, handler.numberInput).
			AddKeybind(0, "", keybindmenu.BindBackspace, handler.backspace)
	} else {
		if showCompleted() {
			menu.AddKeybind('t', "show active tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		} else {
			menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.startDash).
				AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addFromHome).
				AddKeybind('t', "show completed tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		}
		menu.AddKeybind(0, "", keybindmenu.BindNumber, handler.numberInput).
			AddKeybind(0, "", keybindmenu.BindBackspace, handler.backspace)
	}

	menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDown).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUp).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit).
		Apply(home)
}

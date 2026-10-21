package controller

import (
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

func (c *Controller) Home() {
	c.refreshTasks()
	configureCommandLine(c, nil, taskIDInput)
	handler := &keybindHandler{c}

	setWrapper := func() {
		menu := configureHomeKeybinds(handler, c.commandLine.HasInput, c.displayPanel.ShowCompleted)
		c.navPanel.SetText(menu.GenerateText())
		c.commandLine.SetInputCapture(menu.GenerateInputCapture())
	}
	c.commandLine.SetChangedFunc(setWrapper)
	setWrapper()

	c.navPanel.SetTitle("Welcome to Dash!")
}

func configureHomeKeybinds(handler *keybindHandler, hasInput func() bool, showCompleted func() bool) *keybindmenu.KeybindMenu {
	menu := configureKeybindMenu(handler, handler.openTask, taskIDInput)

	if hasInput() {
		menu.AddKeybind('o', "open task", keybindmenu.DefaultBind, handler.openTask).
			AddKeybind('c', "copy task to clipboard", keybindmenu.DefaultBind, handler.copyTask)
		if showCompleted() {
			menu.AddKeybind('u', "reactivate task", keybindmenu.DefaultBind, handler.reactivateTask)
		} else {
			menu.AddKeybind('p', "bump task priority", keybindmenu.DefaultBind, handler.bumpTask)
		}
		menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTask)
	} else {
		if showCompleted() {
			menu.AddKeybind('t', "show active tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		} else {
			menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.startDash).
				AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addFromHome).
				AddKeybind('t', "show completed tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
		}
	}

	return menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDown).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUp).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
}

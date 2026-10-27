package controller

import (
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

type homeMode int

const (
	active homeMode = iota
	completed
	activeWithInput
	completedWithInput
	activeWithNoTasks
	completedWithNoTasks
)

func (c *Controller) Home() {
	c.refreshTasks()
	configureCommandLine(c, nil, taskIDInput)

	setWrapper := func() {
		mode := determineMode(c)
		menu := configureHomeKeybinds(c, mode)
		c.navPanel.SetText(menu.GenerateText())
		if mode == activeWithNoTasks || mode == completedWithNoTasks {
			configureCommandLine(c, menu, keybindInput)
		} else {
			c.commandLine.SetInputCapture(menu.GenerateInputCapture())
		}
	}
	c.commandLine.SetChangedFunc(setWrapper)
	setWrapper()

	c.navPanel.SetTitle("Welcome to Dash!")
}

func determineMode(c *Controller) homeMode {
	isEmpty := c.displayPanel.IsEmpty()
	hasInput := c.commandLine.HasInput()
	showCompleted := c.displayPanel.ShowCompleted()
	if isEmpty && showCompleted {
		return completedWithNoTasks
	} else if isEmpty {
		return activeWithNoTasks
	}

	if hasInput && showCompleted {
		return completedWithInput
	} else if hasInput {
		return activeWithInput
	} else if showCompleted {
		return completed
	} else {
		return active
	}
}

func configureHomeKeybinds(c *Controller, mode homeMode) *keybindmenu.KeybindMenu {
	handler := &keybindHandler{c}
	var menu *keybindmenu.KeybindMenu
	if mode == activeWithNoTasks || mode == completedWithNoTasks {
		menu = configureKeybindMenu(handler, nil, keybindInput)
	} else {
		menu = configureKeybindMenu(handler, handler.openTask, taskIDInput)
	}

	switch mode {
	case active:
		menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.openDash).
			AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addTaskFlow(c.Home)).
			AddKeybind('t', "show completed tasks", keybindmenu.DefaultBind, handler.toggleTaskMode).
			AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDownTable).
			AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUpTable)
	case activeWithInput:
		menu.AddKeybind('o', "open task", keybindmenu.DefaultBind, handler.openTask).
			AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTaskFromHome).
			AddKeybind('p', "bump task priority", keybindmenu.DefaultBind, handler.bumpTaskFromHome)
	case completed:
		menu.AddKeybind('t', "show active tasks", keybindmenu.DefaultBind, handler.toggleTaskMode).
			AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDownTable).
			AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUpTable)
	case completedWithInput:
		menu.AddKeybind('o', "open task", keybindmenu.DefaultBind, handler.openTask).
			AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTaskFromHome).
			AddKeybind('u', "reactivate task", keybindmenu.DefaultBind, handler.reactivateTaskFromHome)
	case activeWithNoTasks:
		menu.AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addTaskFlow(c.Home)).
			AddKeybind('t', "show completed tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
	case completedWithNoTasks:
		menu.AddKeybind('t', "show active tasks", keybindmenu.DefaultBind, handler.toggleTaskMode)
	}
	return menu.AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
}

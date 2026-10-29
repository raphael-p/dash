package controller

import (
	"fmt"
	"strings"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

func (c *Controller) Home() {
	c.refreshTasks()
	configureCommandLine(c, nil, taskIDInput)

	setWrapper := func() {
		menu := configureHomeKeybinds(c)
		c.navPanel.SetText(menu.GenerateText())
		if c.displayPanel.HasTasks() {
			c.commandLine.SetInputCapture(menu.GenerateInputCapture()) // avoid full reset
		} else {
			configureCommandLine(c, menu, keybindInput)
		}
	}
	c.commandLine.SetChangedFunc(setWrapper)
	setWrapper()

	c.navPanel.SetTitle("Welcome to Dash!")
}

func configureHomeKeybinds(c *Controller) *keybindmenu.KeybindMenu {
	taskMode, hasInput, hasTasks := c.displayPanel.GetTaskMode(), c.commandLine.HasInput(), c.displayPanel.HasTasks()

	handler := &keybindHandler{c}
	var menu *keybindmenu.KeybindMenu
	if hasTasks {
		menu = configureKeybindMenu(handler, handler.openTask, taskIDInput)
	} else {
		menu = configureKeybindMenu(handler, nil, keybindInput)
	}

	if hasInput {
		onSuccess := func() { c.refreshTasks(); c.commandLine.Clear() }
		menu.AddKeybind('o', "open task", keybindmenu.DefaultBind, handler.openTask)
		switch taskMode {
		case database.CompletedTasks:
			menu.AddKeybind('u', "reactivate task", keybindmenu.DefaultBind, handler.reactivateTask(nil, onSuccess))
		case database.ActiveTasks:
			menu.AddKeybind('f', "put task first", keybindmenu.DefaultBind, handler.bumpTask(nil, onSuccess))
			menu.AddKeybind('p', "pause task", keybindmenu.DefaultBind, handler.toggleTaskPause(nil, onSuccess))
		case database.PausedTasks:
			menu.AddKeybind('u', "unpause task", keybindmenu.DefaultBind, handler.toggleTaskPause(nil, onSuccess))
		}
		menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTaskFromHome)
	} else {
		if taskMode == database.ActiveTasks && hasTasks {
			menu.AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.openDash)
		}
		if taskMode == database.ActiveTasks {
			menu.AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addTaskFlow(c.Home))
		}

		toggleText := fmt.Sprintf("show %s tasks", strings.ToLower(taskMode.GetNext().ToString()))
		menu.AddKeybind('t', toggleText, keybindmenu.DefaultBind, handler.toggleTaskMode)

		if hasTasks {
			menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDownTable).
				AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUpTable)
		}
	}
	return menu.AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
}

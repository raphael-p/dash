package controller

import (
	"fmt"
	"strings"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func (c *Controller) Home() {
	c.refreshTasks()
	c.commandLine.Clear()

	setWrapper := func() {
		menu, mode := configureHomeKeybinds(c)
		c.navPanel.SetText(menu.Text())
		charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
		c.commandLine.Set(charLimit, placeholder, menu.InputCapture())
	}
	c.commandLine.SetChangedFunc(setWrapper)
	setWrapper()

	c.navPanel.SetTitle("Welcome to Dash!")
}

func configureHomeKeybinds(c *Controller) (*tviewcomponents.KeybindMenu, config.InputMode) {
	taskMode, hasInput, hasTasks := c.displayPanel.GetTaskMode(), c.commandLine.HasInput(), c.displayPanel.HasTasks()

	handler := &keybindHandler{c}
	var menu *tviewcomponents.KeybindMenu
	var mode config.InputMode
	if hasTasks {
		mode = config.TaskIDInput
		menu = handler.newMenu(mode, handler.openTask)
	} else {
		mode = config.KeybindInput
		menu = tviewcomponents.NewKeybindMenu()
	}

	if hasInput {
		onSuccess := func() { c.Home() }
		menu.AddKeybind('o', "open task", tviewcomponents.DefaultBind, handler.openTask)
		switch taskMode {
		case database.CompletedTasks:
			menu.AddKeybind('u', "reactivate task", tviewcomponents.DefaultBind, handler.reactivateTask(nil, onSuccess))
		case database.ActiveTasks:
			menu.AddKeybind('f', "put task first", tviewcomponents.DefaultBind, handler.bumpTask(nil, onSuccess))
			menu.AddKeybind('p', "pause task", tviewcomponents.DefaultBind, handler.toggleTaskPause(nil, onSuccess))
		case database.PausedTasks:
			menu.AddKeybind('u', "unpause task", tviewcomponents.DefaultBind, handler.toggleTaskPause(nil, onSuccess))
		}
		menu.AddKeybind('r', "remove task", tviewcomponents.DefaultBind, handler.removeTaskFromHome)
	} else {
		if taskMode == database.ActiveTasks && hasTasks {
			menu.AddKeybind('d', "dash", tviewcomponents.DefaultBind, handler.openDash)
		}
		if taskMode == database.ActiveTasks {
			menu.AddKeybind('a', "add a new task", tviewcomponents.DefaultBind, handler.addTaskFlow(c.Home))
		}

		toggleText := fmt.Sprintf("show %s tasks", strings.ToLower(taskMode.GetNext().ToString()))
		menu.AddKeybind('t', toggleText, tviewcomponents.DefaultBind, handler.toggleTaskMode)

		if hasTasks {
			menu.AddKeybind('j', "scroll down", tviewcomponents.DefaultBind, handler.scrollDownTable).
				AddKeybind('k', "scroll up", tviewcomponents.DefaultBind, handler.scrollUpTable)
		}
	}
	return menu.AddKeybind('q', "quit", tviewcomponents.DefaultBind, handler.quit), mode
}

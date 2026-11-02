package controller

import (
	"fmt"
	"strings"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func (c *Controller) Home() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
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

	var menu *tviewcomponents.KeybindMenu
	var mode config.InputMode
	if hasTasks {
		mode = config.TaskIDInput
		menu = newMenuForInput(c, mode, openTask(c))
	} else {
		mode = config.KeybindInput
		menu = tviewcomponents.NewKeybindMenu()
	}

	if hasInput {
		onSuccess := func() { c.Home() }
		menu.AddKeybind('o', "open task", tviewcomponents.DefaultBind, openTask(c))
		switch taskMode {
		case database.CompletedTasks:
			menu.AddKeybind('u', "reactivate task", tviewcomponents.DefaultBind, reactivateTask(c, nil, onSuccess))
		case database.ActiveTasks:
			menu.AddKeybind('f', "put task first", tviewcomponents.DefaultBind, bumpTask(c, nil, onSuccess))
			menu.AddKeybind('p', "pause task", tviewcomponents.DefaultBind, toggleTaskPause(c, nil, onSuccess))
		case database.PausedTasks:
			menu.AddKeybind('u', "unpause task", tviewcomponents.DefaultBind, toggleTaskPause(c, nil, onSuccess))
		}
		menu.AddKeybind('r', "remove task", tviewcomponents.DefaultBind, removeTaskFromHome(c))
	} else {
		if taskMode == database.ActiveTasks && hasTasks {
			menu.AddKeybind('d', "dash", tviewcomponents.DefaultBind, c.openDash)
		}
		if taskMode == database.ActiveTasks {
			menu.AddKeybind('a', "add task", tviewcomponents.DefaultBind, func() { c.addTaskName(c.Home) })
		}

		toggleText := fmt.Sprintf("show %s", strings.ToLower(taskMode.GetNext().ToString()))
		menu.AddKeybind('t', toggleText, tviewcomponents.DefaultBind, func() { c.displayPanel.ToggleTaskMode(); c.Home() })

		if hasTasks {
			menu.AddKeybind('j', "scroll down", tviewcomponents.DefaultBind, scrollDownHome(c)).
				AddKeybind('k', "scroll up", tviewcomponents.DefaultBind, scrollUpHome(c))
		}
	}
	return menu.AddKeybind('q', "quit", tviewcomponents.DefaultBind, c.stop), mode
}

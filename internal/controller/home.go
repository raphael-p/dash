package controller

import (
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
		var onEnter, onScrollDown, onScrollUp tviewcomponents.KeybindHandler
		if hasInput {
			onEnter = openTask(c)
		} else {
			onScrollDown, onScrollUp = scrollDownHome(c), scrollUpHome(c)
		}
		menu = newInputKeybindMenu(c, mode, onEnter, nil, onScrollDown, onScrollUp)
	} else {
		mode = config.KeybindInput
		menu = newInputKeybindMenu(c, mode, nil, nil, nil, nil)
	}

	if hasInput {
		onSuccess := func() { c.Home() }
		menu.AddRuneKeybind('o', "open task", openTask(c))
		switch taskMode {
		case database.CompletedTasks:
			menu.AddRuneKeybind('u', "reactivate task", reactivateTask(c, nil, onSuccess))
		case database.ActiveTasks:
			menu.AddRuneKeybind('f', "put task first", bumpTask(c, nil, onSuccess))
			menu.AddRuneKeybind('p', "pause task", toggleTaskPause(c, nil, onSuccess))
		case database.PausedTasks:
			menu.AddRuneKeybind('u', "unpause task", toggleTaskPause(c, nil, onSuccess))
		}
		menu.AddRuneKeybind('r', "remove task", removeTaskFromHome(c))
	} else {
		if taskMode == database.ActiveTasks && hasTasks {
			menu.AddRuneKeybind('d', "dash", c.openDash)
		}
		if taskMode == database.ActiveTasks {
			menu.AddRuneKeybind('a', "add task", func() { c.addTaskName(c.Home) })
		}

		toggleText := "show " + strings.ToLower(taskMode.GetNext().ToString())
		menu.AddPrimaryKeybind("TAB", toggleText, tviewcomponents.BindTab, func() { c.displayPanel.ToggleTaskMode(); c.Home() })
	}
	return menu, mode
}

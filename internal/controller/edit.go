package controller

import (
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func (c *Controller) editTask(task *database.Task, back func()) {
	var onScrollDown, onScrollUp tviewcomponents.KeybindHandler
	if strutil.LineCount(task.Description) > 2 { // heuristic to reduce visual clutter
		onScrollDown, onScrollUp = scrollDownEdit(c), scrollUpEdit(c)
	}
	menu := newInputKeybindMenu(c, config.KeybindInput, nil, func() { back(); c.infoPanel.Clear() }, onScrollDown, onScrollUp)

	if task.CompletedAt.Valid {
		menu.AddRuneKeybind('u', "reactivate", reactivateTask(c, task, func() {}))
	} else if task.PausedAt.Valid {
		menu.AddRuneKeybind('u', "unpause", toggleTaskPause(c, task, func() {}))
	}

	if !task.CompletedAt.Valid {
		refresh := func() {
			if task, ok := database.Service().GetTask(task.ID); ok {
				c.displayPanel.ShowTask(task)
				c.editTask(&task, back)
			}
		}
		menu.AddRuneKeybind('n', "edit name", func() { c.editTaskName(task, refresh) }).
			AddRuneKeybind('d', "edit description", func() { c.editTaskDescription(task, refresh) })
	}

	if !task.CompletedAt.Valid && !task.PausedAt.Valid {
		menu.AddRuneKeybind('f', "put first", bumpTask(c, task, func() {})).
			AddRuneKeybind('p', "pause", toggleTaskPause(c, task, func() {}))

	}

	menu.AddRuneKeybind('c', "copy description", copy(c, task.Description)).
		AddRuneKeybind('r', "remove", removeTaskFromEdit(c, task, back))

	c.commandLine.Clear().Hide(menu.InputCapture())
	c.navPanel.SetTitle("Edit Task").SetText(menu.Text())
}

func (c *Controller) editTaskName(task *database.Task, refresh func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newInputKeybindMenu(c, mode, submitName(c, task, refresh), refresh, nil, nil)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Name)
	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.Text())
}

func (c *Controller) editTaskDescription(task *database.Task, refresh func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newInputKeybindMenu(c, mode, submitDescription(c, task, refresh), refresh, nil, nil)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Description)
	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.Text())
}

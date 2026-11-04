package controller

import (
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func (c *Controller) editTask(task *database.Task, back func()) {
	menu := newInputKeybindMenu(c, config.KeybindInput, nil, func() { back(); c.infoPanel.Clear() })

	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate", tviewcomponents.DefaultBind, reactivateTask(c, task, func() {}))
	} else if task.PausedAt.Valid {
		menu.AddKeybind('u', "unpause", tviewcomponents.DefaultBind, toggleTaskPause(c, task, func() {}))
	}

	if !task.CompletedAt.Valid {
		refresh := func() {
			if task, ok := database.Service().GetTask(task.ID); ok {
				c.displayPanel.ShowTask(task)
				c.editTask(&task, back)
			}
		}
		menu.AddKeybind('n', "edit name", tviewcomponents.DefaultBind, func() { c.editTaskName(task, refresh) }).
			AddKeybind('d', "edit description", tviewcomponents.DefaultBind, func() { c.editTaskDescription(task, refresh) })
	}

	if !task.CompletedAt.Valid && !task.PausedAt.Valid {
		menu.AddKeybind('f', "put first", tviewcomponents.DefaultBind, bumpTask(c, task, func() {})).
			AddKeybind('p', "pause", tviewcomponents.DefaultBind, toggleTaskPause(c, task, func() {}))

	}

	menu.AddKeybind('c', "copy description", tviewcomponents.DefaultBind, copy(c, task.Description)).
		AddKeybind('r', "remove", tviewcomponents.DefaultBind, removeTaskFromEdit(c, task, back))

	if strutil.LineCount(task.Description) > 2 { // heuristic to reduce visual clutter
		menu.AddKeybind('j', "scroll down", tviewcomponents.DefaultBind, scrollDownEdit(c)).
			AddKeybind('k', "scroll up", tviewcomponents.DefaultBind, scrollUpEdit(c))
	}

	c.commandLine.Clear().Hide(menu.InputCapture())
	c.navPanel.SetTitle("Edit Task").SetText(menu.Text())
}

func (c *Controller) editTaskName(task *database.Task, refresh func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newInputKeybindMenu(c, mode, submitName(c, task, refresh), refresh)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Name)
	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.Text())
}

func (c *Controller) editTaskDescription(task *database.Task, refresh func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newInputKeybindMenu(c, mode, submitDescription(c, task, refresh), refresh)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Description)
	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.Text())
}

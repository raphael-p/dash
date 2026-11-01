package controller

import (
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func (c *Controller) editTask(task *database.Task, back func()) {
	mode := config.KeybindInput
	menu := tviewcomponents.NewKeybindMenu()

	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate", tviewcomponents.DefaultBind, reactivateTask(c, task, func() {}))
	} else if task.PausedAt.Valid {
		menu.AddKeybind('u', "unpause", tviewcomponents.DefaultBind, toggleTaskPause(c, task, func() {}))
	}

	if !task.CompletedAt.Valid {
		onSuccess := func() {
			if task, ok := database.Service().GetTask(task.ID); ok {
				c.displayPanel.ShowTask(task)
				c.editTask(&task, back)
			}
		}
		menu.AddKeybind('n', "change name", tviewcomponents.DefaultBind, func() { c.editTaskName(task, onSuccess) }).
			AddKeybind('d', "change description", tviewcomponents.DefaultBind, func() { c.editTaskDescription(task, onSuccess) })
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
	menu.AddKeybind('b', "back", tviewcomponents.DefaultBind, func() { back(); c.infoPanel.Clear() }).
		AddKeybind('q', "quit", tviewcomponents.DefaultBind, c.stop)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Edit Task").SetText(menu.Text())
}

func (c *Controller) editTaskName(task *database.Task, onSuccess func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newMenuForInput(c, mode, submitName(c, task, onSuccess))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Name)
	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.Text())
}

func (c *Controller) editTaskDescription(task *database.Task, onSuccess func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newMenuForInput(c, mode, submitDescription(c, task, onSuccess))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Description)
	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.Text())
}

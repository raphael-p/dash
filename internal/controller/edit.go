package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
	"github.com/rivo/tview"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newMenuForInput(c, mode, addTaskFlowSubmitName(c, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (name)").SetText(fmt.Sprintf("name: ?\n\ndescription: ?\n\n%s", menu.Text()))
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newMenuForInput(c, mode, addTaskFlowSubmitDescription(c, name, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (description)").
		SetText(fmt.Sprintf("name: %s\n\ndescription: ?\n\n%s", tview.Escape(name), menu.Text()))
}

func (c *Controller) addTaskInsertionOrder(name, description string, back func()) {
	c.infoPanel.Clear()

	mode := config.KeybindInput
	submit := func(insertionOrder database.TaskInsertionOrder) func() {
		return func() { database.Service().AddTask(name, description, insertionOrder, back) }
	}
	menu := tviewcomponents.NewKeybindMenu().
		AddKeybind('f', "insert first", tviewcomponents.DefaultBind, submit(database.InsertFirst)).
		AddKeybind('n', "insert next", tviewcomponents.DefaultBind, submit(database.InsertNext)).
		AddKeybind('l', "insert last", tviewcomponents.DefaultBind, submit(database.InsertLast))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.
		SetTitle("Add New Task (insert)").
		SetText(fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", tview.Escape(name), tview.Escape(description), menu.Text()))
}

func (c *Controller) editTaskFlow(task *database.Task, back func()) {
	mode := config.KeybindInput
	menu := tviewcomponents.NewKeybindMenu()

	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate", tviewcomponents.DefaultBind, reactivateTask(c, task, func() {}))
	} else if task.PausedAt.Valid {
		menu.AddKeybind('u', "unpause", tviewcomponents.DefaultBind, toggleTaskPause(c, task, func() {}))
	}

	if !task.CompletedAt.Valid {
		menu.AddKeybind('n', "change name", tviewcomponents.DefaultBind, func() { c.editTaskName(task, back) }).
			AddKeybind('d', "change description", tviewcomponents.DefaultBind, func() { c.editTaskDescription(task, back) })
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

func (c *Controller) editTaskName(task *database.Task, back func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newMenuForInput(c, mode, editTaskFlowSubmitName(c, task, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Name)
	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.Text())
}

func (c *Controller) editTaskDescription(task *database.Task, back func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newMenuForInput(c, mode, editTaskFlowSubmitDescription(c, task, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Description)
	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.Text())
}

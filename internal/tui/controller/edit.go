package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	mode := config.NameInput
	menu := handler.newMenu(mode, handler.addTaskFlowSubmitName(back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (name)").SetText(fmt.Sprintf("name: ?\n\ndescription: ?\n\n%s", menu.Text()))
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	mode := config.DescriptionInput
	menu := handler.newMenu(mode, handler.addTaskFlowSubmitDescription(name, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (description)").
		SetText(fmt.Sprintf("name: %s\n\ndescription: ?\n\n%s", tview.Escape(name), menu.Text()))
}

func (c *Controller) addTaskInsertionOrder(name, description string, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	mode := config.KeybindInput
	menu := keybindmenu.New().
		AddKeybind('f', "insert first", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertFirst, back)).
		AddKeybind('n', "insert next", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertNext, back)).
		AddKeybind('l', "insert last", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertLast, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.
		SetTitle("Add New Task (insert)").
		SetText(fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", tview.Escape(name), tview.Escape(description), menu.Text()))
}

func (c *Controller) editTaskFlow(task *database.Task, back func()) {
	handler := &keybindHandler{c}
	mode := config.KeybindInput
	menu := keybindmenu.New()

	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate", keybindmenu.DefaultBind, handler.reactivateTask(task, func() {}))
	} else if task.PausedAt.Valid {
		menu.AddKeybind('u', "unpause", keybindmenu.DefaultBind, handler.toggleTaskPause(task, func() {}))
	}

	if !task.CompletedAt.Valid {
		menu.AddKeybind('n', "change name", keybindmenu.DefaultBind, handler.editTaskFlowName(task, back)).
			AddKeybind('d', "change description", keybindmenu.DefaultBind, handler.editTaskFlowDescription(task, back))
	}

	if !task.CompletedAt.Valid && !task.PausedAt.Valid {
		menu.AddKeybind('f', "put first", keybindmenu.DefaultBind, handler.bumpTask(task, func() {})).
			AddKeybind('p', "pause", keybindmenu.DefaultBind, handler.toggleTaskPause(task, func() {}))

	}

	menu.AddKeybind('c', "copy description", keybindmenu.DefaultBind, handler.copy(task.Description)).
		AddKeybind('r', "remove", keybindmenu.DefaultBind, handler.removeTaskFromEdit(task, back))

	if strutil.LineCount(task.Description) > 2 { // heuristic to reduce visual clutter
		menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDownTask).
			AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUpTask)
	}
	menu.AddKeybind('b', "back", keybindmenu.DefaultBind, func(_ rune) { back(); c.infoPanel.Clear() }).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Edit Task").SetText(menu.Text())
}

func (c *Controller) editTaskName(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	mode := config.NameInput
	menu := handler.newMenu(mode, handler.editTaskFlowSubmitName(task, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Name)
	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.Text())
}

func (c *Controller) editTaskDescription(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	mode := config.DescriptionInput
	menu := handler.newMenu(mode, handler.editTaskFlowSubmitDescription(task, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.commandLine.SetInput(task.Description)
	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.Text())
}

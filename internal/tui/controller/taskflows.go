package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/stringpad"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.addTaskFlowSubmitName(back), nameInput)
	configureCommandLine(c, menu, nameInput)

	c.navPanel.SetTitle("Add New Task (name)").SetText(fmt.Sprintf("name: ?\n\ndescription: ?\n\n%s", menu.GenerateText()))
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(&keybindHandler{c}, handler.addTaskFlowSubmitDescription(name, back), descriptionInput)
	configureCommandLine(c, menu, descriptionInput)

	c.navPanel.SetTitle("Add New Task (description)").SetText(fmt.Sprintf("name: %s\n\ndescription: ?\n\n%s", name, menu.GenerateText()))
}

func (c *Controller) addTaskInsertionOrder(name, description string, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(&keybindHandler{c}, nil, keybindInput).
		AddKeybind('f', "insert first", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertFirst, back)).
		AddKeybind('n', "insert next", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertNext, back)).
		AddKeybind('l', "insert last", keybindmenu.DefaultBind, handler.addTaskFlowSubmit(name, description, InsertLast, back))
	configureCommandLine(c, menu, keybindInput)

	c.navPanel.SetTitle("Add New Task (insert)").SetText(fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", name, description, menu.GenerateText()))
}

func (c *Controller) editTaskFlow(task *database.Task, back func()) {
	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, nil, keybindInput)

	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate", keybindmenu.DefaultBind, handler.reactivateTaskFromEdit(task.ID))
	} else if task.PausedAt.Valid {
		menu.AddKeybind('u', "unpause", keybindmenu.DefaultBind, handler.unpauseTaskFromEdit(task.ID))
	}

	if !task.CompletedAt.Valid {
		menu.AddKeybind('n', "change name", keybindmenu.DefaultBind, handler.editTaskFlowName(task, back)).
			AddKeybind('d', "change description", keybindmenu.DefaultBind, handler.editTaskFlowDescription(task, back))
	}

	if !task.CompletedAt.Valid && !task.PausedAt.Valid {
		menu.AddKeybind('f', "put first", keybindmenu.DefaultBind, handler.bumpTaskFromEdit(task.ID)).
			AddKeybind('p', "pause", keybindmenu.DefaultBind, handler.pauseTaskFromEdit(task.ID))

	}

	menu.AddKeybind('c', "copy description", keybindmenu.DefaultBind, handler.copy(task.Description)).
		AddKeybind('r', "remove", keybindmenu.DefaultBind, handler.removeTaskFromEdit(task, back))

	if stringpad.LineCount(task.Description) > 2 { // heuristic to reduce visual clutter
		menu.AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDownTask).
			AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUpTask)
	}
	menu.AddKeybind('b', "back", keybindmenu.DefaultBind, func(_ rune) { back() }).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
	configureCommandLine(c, menu, keybindInput)

	c.navPanel.SetTitle("Edit Task").SetText(menu.GenerateText())
}

func (c *Controller) editTaskName(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.editTaskFlowSubmitName(task, back), nameInput)
	configureCommandLine(c, menu, nameInput)
	c.commandLine.SetInput(task.Name)

	c.navPanel.SetTitle("Edit Task (name)").SetText(menu.GenerateText())
}

func (c *Controller) editTaskDescription(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.editTaskFlowSubmitDescription(task, back), descriptionInput)
	configureCommandLine(c, menu, descriptionInput)
	c.commandLine.SetInput(task.Description)

	c.navPanel.SetTitle("Edit Task (description)").SetText(menu.GenerateText())
}

package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.addTaskFlowSubmitName(back), nameInput)
	configureCommandLine(c, menu, nameInput)

	c.navPanel.SetTitle("Add New Task (name)").SetText("name: ?\n\ndescription: ?")
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(&keybindHandler{c}, handler.addTaskFlowSubmitDescription(name, back), descriptionInput)
	configureCommandLine(c, menu, descriptionInput)

	c.navPanel.SetTitle("Add New Task (description)").SetText(fmt.Sprintf("name: %s\n\ndescription: ?", name))
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
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, nil, keybindInput)
	if task.CompletedAt.Valid {
		menu.AddKeybind('u', "reactivate task", keybindmenu.DefaultBind, handler.reactivateTaskFromEdit(task.ID))
	} else {
		menu.AddKeybind('n', "change name", keybindmenu.DefaultBind, handler.editTaskFlowName(task, back)).
			AddKeybind('d', "change description", keybindmenu.DefaultBind, handler.editTaskFlowDescription(task, back)).
			AddKeybind('p', "bump task priority", keybindmenu.DefaultBind, handler.bumpTaskFromEdit(task.ID))
	}
	menu.AddKeybind('r', "remove task", keybindmenu.DefaultBind, handler.removeTaskFromEdit(task.ID, back)).
		AddKeybind('b', "back", keybindmenu.DefaultBind, func(_ rune) { back() }).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
	configureCommandLine(c, menu, keybindInput)

	c.navPanel.SetTitle("Edit Task").SetText(menu.GenerateText())
}

func (c *Controller) editTaskName(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.editTaskFlowSubmitName(task, back), nameInput)
	configureCommandLine(c, menu, nameInput)

	c.navPanel.SetTitle("Edit Task (name)").SetText("")
}

func (c *Controller) editTaskDescription(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, handler.editTaskFlowSubmitDescription(task, back), descriptionInput)
	configureCommandLine(c, menu, descriptionInput)

	c.navPanel.SetTitle("Edit Task (description)").SetText("")
}

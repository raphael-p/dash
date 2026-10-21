package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	onEnter := func(_ rune) {
		if name, ok := getInputOrWarn(c.commandLine.GetInput(), "name", c.infoPanel.Warn); ok {
			c.addTaskDescription(back, name)
		}
	}
	menu := configureKeybindMenu(&keybindHandler{c}, onEnter, nameInput)
	configureCommandLine(c, menu, nameInput)

	c.navPanel.SetTitle("Add New Task (name)").SetText("name: ?\n\ndescription: ?")
}

func (c *Controller) addTaskDescription(back func(), name string) {
	c.infoPanel.Clear()

	onEnter := func(_ rune) {
		if description, ok := getInputOrWarn(c.commandLine.GetInput(), "description", c.infoPanel.Warn); ok {
			c.addTaskInsertionOrder(back, name, description)
		}
	}
	menu := configureKeybindMenu(&keybindHandler{c}, onEnter, descriptionInput)
	configureCommandLine(c, menu, descriptionInput)

	c.navPanel.SetTitle("Add New Task (description)").SetText(fmt.Sprintf("name: %s\n\ndescription: ?", name))
}

func (c *Controller) addTaskInsertionOrder(back func(), name, description string) {
	c.infoPanel.Clear()

	submit := func(insertionOrder taskInsertionOrder) func(rune) {
		return func(_ rune) {
			ok := c.addTask(name, description, insertionOrder)
			if ok {
				back()
			}
		}
	}
	menu := configureKeybindMenu(&keybindHandler{c}, nil, keybindInput).
		AddKeybind('f', "insert first", keybindmenu.DefaultBind, submit(InsertFirst)).
		AddKeybind('n', "insert next", keybindmenu.DefaultBind, submit(InsertNext)).
		AddKeybind('l', "insert last", keybindmenu.DefaultBind, submit(InsertLast))
	configureCommandLine(c, menu, keybindInput)

	c.navPanel.SetTitle("Add New Task (insert)").SetText(fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", name, description, menu.GenerateText()))
}

func (c *Controller) editTaskForm(task *database.Task, back func()) {
	c.infoPanel.Clear()

	handler := &keybindHandler{c}
	menu := configureKeybindMenu(handler, nil, keybindInput).
		AddKeybind('n', "change name", keybindmenu.DefaultBind, func(_ rune) { c.editTaskName(task, back) }).
		AddKeybind('d', "change description", keybindmenu.DefaultBind, func(_ rune) { c.editTaskDescription(task, back) }).
		AddKeybind('p', "bump task priority", keybindmenu.DefaultBind, func(_ rune) { c.bumpTask(task.ID) }).
		AddKeybind('b', "back", keybindmenu.DefaultBind, func(_ rune) { back() }).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit)
	configureCommandLine(c, menu, keybindInput)

	c.navPanel.SetTitle("Edit Task").SetText(menu.GenerateText())
}

func (c *Controller) editTaskName(task *database.Task, back func()) {
	c.infoPanel.Clear()

	onEnter := func(_ rune) {
		if name, ok := getInputOrWarn(c.commandLine.GetInput(), "name", c.infoPanel.Warn); ok {
			if c.editTask(task, name, "") {
				c.openTask(task.ID, back)
			}
		}
	}
	menu := configureKeybindMenu(&keybindHandler{c}, onEnter, nameInput)
	configureCommandLine(c, menu, nameInput)

	c.navPanel.SetTitle("Edit Task (name)").SetText("")
}

func (c *Controller) editTaskDescription(task *database.Task, back func()) {
	c.infoPanel.Clear()

	onEnter := func(_ rune) {
		if description, ok := getInputOrWarn(c.commandLine.GetInput(), "description", c.infoPanel.Warn); ok {
			if c.editTask(task, "", description) {
				c.openTask(task.ID, back)
			}
		}
	}
	menu := configureKeybindMenu(&keybindHandler{c}, onEnter, descriptionInput)
	configureCommandLine(c, menu, descriptionInput)

	c.navPanel.SetTitle("Edit Task (description)").SetText("")
}

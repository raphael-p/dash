package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/rivo/tview"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newInputKeybindMenu(c, mode, submitNameCandidate(c, back), back)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (name)").SetText(fmt.Sprintf("name: ?\n\ndescription: ?\n\n%s", menu.Text()))
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newInputKeybindMenu(c, mode, submitDescriptionCandidate(c, name, back), back)

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (description)").
		SetText(fmt.Sprintf("name: %s\n\ndescription: ?\n\n%s", tview.Escape(name), menu.Text()))
}

func (c *Controller) addTaskInsertionOrder(name, description string, back func()) {
	c.infoPanel.Clear()

	submit := func(insertionOrder database.TaskInsertionOrder) func() {
		return func() { database.Service().AddTask(name, description, insertionOrder, back) }
	}
	menu := newInputKeybindMenu(c, config.KeybindInput, nil, back).
		AddRuneKeybind('f', "insert first", submit(database.InsertFirst)).
		AddRuneKeybind('n', "insert next", submit(database.InsertNext)).
		AddRuneKeybind('l', "insert last", submit(database.InsertLast))

	c.commandLine.Clear().Hide(menu.InputCapture())
	text := fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", tview.Escape(name), tview.Escape(description), menu.Text())
	c.navPanel.SetTitle("Add New Task (insert)").SetText(text)
}

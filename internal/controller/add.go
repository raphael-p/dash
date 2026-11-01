package controller

import (
	"fmt"

	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
	"github.com/rivo/tview"
)

func (c *Controller) addTaskName(back func()) {
	c.infoPanel.Clear()

	mode := config.NameInput
	menu := newMenuForInput(c, mode, submitNameCandidate(c, back))

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
	c.navPanel.SetTitle("Add New Task (name)").SetText(fmt.Sprintf("name: ?\n\ndescription: ?\n\n%s", menu.Text()))
}

func (c *Controller) addTaskDescription(name string, back func()) {
	c.infoPanel.Clear()

	mode := config.DescriptionInput
	menu := newMenuForInput(c, mode, submitDescriptionCandidate(c, name, back))

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
	text := fmt.Sprintf("name: %s\n\ndescription: %s\n\n%s", tview.Escape(name), tview.Escape(description), menu.Text())
	c.navPanel.SetTitle("Add New Task (insert)").SetText(text)
}

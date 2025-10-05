package controller

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/tui/panels"
	"github.com/rivo/tview"
)

type Controller struct {
	app          *tview.Application
	inputPanel   *panels.InputPanel
	infoPanel    *panels.InfoPanel
	displayPanel *panels.DisplayPanel
}

func NewController(
	app *tview.Application,
	inputPanel *panels.InputPanel,
	infoPanel *panels.InfoPanel,
	displayPanel *panels.DisplayPanel,
) *Controller {
	return &Controller{app, inputPanel, infoPanel, displayPanel}
}

func (c *Controller) MainMenu() {
	c.infoPanel.Clear()
	c.refreshTasks()

	navigationList := tview.NewList()
	navigationList.AddItem("Open Task", "", 'o', func() { c.openTaskForm() })
	navigationList.AddItem("Add Task", "", 'a', func() { c.addTaskForm() })
	navigationList.AddItem("Remove Task", "", 'r', func() { c.removeTaskForm() })
	navigationList.AddItem("Scroll Down", "", 'j', func() {
		err := c.displayPanel.ScrollDown()
		if err != nil {
			c.infoPanel.Error(fmt.Errorf("failed to scroll down: %s", err))
		}
	})
	navigationList.AddItem("Scroll Up", "", 'k', func() {
		err := c.displayPanel.ScrollUp()
		if err != nil {
			c.infoPanel.Error(fmt.Errorf("failed to scroll up: %s", err))
		}
	})
	navigationList.AddItem("Quit", "", 'q', func() { c.app.Stop() })
	c.inputPanel.Set("Manage Your Tasks", navigationList)

	// handle typing task number directly to open a task
	navigationList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		e := event.Key()
		r := event.Rune()
		switch {
		case e == tcell.KeyEscape || e == tcell.KeyTab:
			c.infoPanel.Clear()
		case e == tcell.KeyEnter && c.infoPanel.GetInput() != "":
			c.submitOpenTaskString(c.infoPanel.GetInput())
			return nil
		case r >= '0' && r <= '9':
			c.infoPanel.AppendInput(string(r))
			return nil
		}
		return event
	})
}

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
}

package controller

import (
	"fmt"

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
	c.refreshTasks()

	navigationList := tview.NewList()
	navigationList.AddItem("View Task", "", '0', func() { c.viewTaskForm() })
	navigationList.AddItem("Add Task", "", '1', func() { c.addTaskForm() })
	navigationList.AddItem("Delete Task", "", '2', func() { c.deleteTaskForm() })
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
}

func (c *Controller) refreshTasks() {
	err := c.displayPanel.GetCurrentPage()
	if err != nil {
		c.infoPanel.Error(err)
	}
}

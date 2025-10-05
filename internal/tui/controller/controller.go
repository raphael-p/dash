package controller

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/database"
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
		case e == tcell.KeyEscape, e == tcell.KeyTab:
			c.infoPanel.Clear()
		case e == tcell.KeyBackspace, e == tcell.KeyBackspace2:
			currentInput := c.infoPanel.GetInput()
			if currentInput != "" {
				c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
			}
		case e == tcell.KeyEnter && c.infoPanel.GetInput() != "":
			c.submitOpenTaskString(c.infoPanel.GetInput())
			return nil
		case r == 'b':
			taskID, err := extractIDFromString(c.infoPanel.GetInput())
			if err != nil {
				c.infoPanel.Warn(fmt.Sprint("your input is invalid: ", err))
				return nil
			}

			bumped, err := database.BumpTask(int64(taskID))
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to bump task priority: %s", err))
				return nil
			}

			if bumped {
				c.infoPanel.Info(fmt.Sprintf("Bumped priority of task [%d]", taskID))
				c.refreshTasks()
			} else {
				c.infoPanel.Warn(fmt.Sprintf("Task [%d] does not exist, noop.", taskID))
			}
		case r >= '0' && r <= '9':
			c.infoPanel.SetInput(c.infoPanel.GetInput() + string(r))
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

package controller

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.infoPanel.Clear()
	c.refreshTasks()

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)
	home.SetText(fmt.Sprintf(
		`([%[1]s::b]d[-:-:-]) dash
([%[1]s::b]o[-:-:-]) open a task
([%[1]s::b]b[-:-:-]) bump the priority of a task
([%[1]s::b]r[-:-:-]) remove a task
([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]j[-:-:-]) scroll down
([%[1]s::b]k[-:-:-]) scroll up
([%[1]s::b]q[-:-:-]) quit
[::d]tip: you can type the task ID before invoking the open (o), bump (b), or remove (r) commands[::-]`,
		tcell.ColorLimeGreen))
	c.inputPanel.Set("Welcome to Dash!", home)

	home.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		e := event.Key()
		r := event.Rune()
		switch {
		// ADD
		case r == 'a':
			c.addTaskForm()
		// BACKSPACE
		case e == tcell.KeyBackspace, e == tcell.KeyBackspace2:
			currentInput := c.infoPanel.GetInput()
			if currentInput != "" {
				c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
			}
		// BUMP
		case r == 'b':
			if c.infoPanel.GetInput() != "" {
				c.bumpTask(c.infoPanel.GetInput())
			} else {
				c.bumpTaskForm()
			}
		// OPEN
		case r == 'o', e == tcell.KeyEnter:
			if c.infoPanel.GetInput() != "" {
				c.openTask(c.infoPanel.GetInput())
			} else {
				c.openTaskForm()
			}
		// QUIT
		case r == 'q':
			c.app.Stop()
		// REMOVE
		case r == 'r':
			if c.infoPanel.GetInput() != "" {
				c.removeTask(c.infoPanel.GetInput())
			} else {
				c.removeTaskForm()
			}
		// SCROLL TEXT
		case e == tcell.KeyDown, e == tcell.KeyUp:
			return event
		// SCROLL DOWN
		case r == 'j':
			err := c.displayPanel.ScrollDown()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to scroll down: %s", err))
			}
		// SCROLL UP
		case r == 'k':
			err := c.displayPanel.ScrollUp()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to scroll up: %s", err))
			}
		// START DASH
		case r == 'd':
			_, err := c.displayPanel.ShowTopTask()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to start timer: %s", err))
			} else {
				c.startDash()
			}
		// TYPE TASK ID
		case r >= '0' && r <= '9':
			c.infoPanel.SetInput(c.infoPanel.GetInput() + string(r))
		default:
			message := "invalid command: "
			if r == 0 || (string(r) != " " && unicode.IsSpace(r)) {
				c.infoPanel.Warn(message + event.Name())
			} else {
				c.infoPanel.Warn(message + "'" + string(r) + "'")
			}
		}
		return nil
	})
}

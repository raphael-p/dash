package controller

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/tui/components/countdowntimer"
)

func (c *Controller) startDash(quit func()) {
	task, err := c.displayPanel.ShowTopTask()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to start dash: %s", err))
		return
	}

	timer := countdowntimer.Instance()

	timer.SetDescription(fmt.Sprintf(
		`([%[1]s::b]c[-:-:-]) mark task as completed
([%[1]s::b]e[-:-:-]) edit task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]b[-:-:-]) back`,
		tcell.ColorLimeGreen))

	timer.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		r := event.Rune()
		back := func() {
			c.startDash(quit)
		}

		switch r {
		// ADD
		case 'a':
			c.addTaskForm(back, quit)
		// BACK
		case 'b':
			c.Home()
		// EDIT
		case 'e':
			c.editTaskForm(task, back, quit)
		// RESET TIMER
		case 'r':
			timer.Reset(c.app)
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

	c.inputPanel.Set("Dash", timer.Layout)

	timer.Start(c.app)
}

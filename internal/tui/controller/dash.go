package controller

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/tui/components/countdowntimer"
)

func (c *Controller) startDash() {
	timer := countdowntimer.Instance()

	timer.SetDescription(fmt.Sprintf(
		`([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]c[-:-:-]) mark task as completed
([%[1]s::b]b[-:-:-]) back`,
		tcell.ColorLimeGreen))

	timer.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		r := event.Rune()
		switch r {
		// ADD
		case 'a':
			c.addTaskForm()
		// BACK
		case 'b':
			c.Home()
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

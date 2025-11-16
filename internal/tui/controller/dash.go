package controller

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/components/countdowntimer"
)

var lastTask database.Task

func (c *Controller) startDash(quit func()) {
	task, err := c.displayPanel.ShowTopTask()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to start dash: %s", err))
		return
	}
	if lastTask.Id == 0 {
		lastTask = task
	}

	timer := countdowntimer.Instance("lock in.", "dash complete. restart when ready.", c.dashDuration)

	timer.SetDescription(fmt.Sprintf(
		`([%[1]s::b]d[-:-:-]) mark task as done, ([%[1]s::b]u[-:-:-]) to undo
([%[1]s::b]e[-:-:-]) edit task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]b[-:-:-]) back`,
		tcell.ColorLimeGreen))

	timer.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		r := event.Rune()
		refresh := func() {
			c.startDash(quit)
		}

		switch r {
		// ADD
		case 'a':
			c.addTaskForm(refresh, quit)
		// BACK
		case 'b':
			c.Home()
		// DONE
		case 'd':
			err := task.MarkAsDone()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to mark task as done: %s", err))
			} else {
				lastTask = task
				refresh()
			}
		// UNDONE
		case 'u':
			if lastTask.Id == task.Id {
				break // noop
			}
			err := lastTask.UndoMarkAsDone()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to undo mark task as done: %s", err))
			} else {
				refresh()
			}
		// EDIT
		case 'e':
			c.editTaskForm(task, refresh, quit)
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

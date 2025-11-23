package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/components/countdowntimer"
)

var lastTask database.Task

func trackTime(c *Controller, task database.Task, dashDuration time.Duration) {
	newTimeSpent := time.Duration(task.TimeSpentSeconds.Int16)*time.Second + dashDuration
	task.TimeSpentSeconds = sql.NullInt16{
		Int16: int16(newTimeSpent.Seconds()),
		Valid: true,
	}
	ok, err := task.Update()
	if !ok || err != nil {
		if err == nil {
			err = errors.New("noop")
		}
		c.infoPanel.Error(fmt.Errorf("failed to update time spent on task %d: %s", task.Id, err))
	}

	var timeSpentMessage string
	if minutesSpent := newTimeSpent.Minutes(); minutesSpent > 1 {
		timeSpentMessage = fmt.Sprintf(": %d minutes spent in total", int(minutesSpent))
	}
	c.infoPanel.Info(fmt.Sprintf("time tracking on task %d updated%s", task.Id, timeSpentMessage))
}

func setDashCommands(c *Controller, quit func(), task database.Task, resetTimer func()) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		e := event.Key()
		if e == tcell.KeyDown || e == tcell.KeyUp {
			return event
		}

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
			task.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
			_, err := task.Update()
			if err != nil {
				task.CompletedAt = sql.NullTime{Valid: false}
				c.infoPanel.Error(fmt.Errorf("failed to mark task %d as done: %s", task.Id, err))
			} else {
				lastTask = task
				refresh()
			}
		// UNDONE
		case 'u':
			if lastTask.Id == task.Id {
				break // noop
			}
			lastTask.CompletedAt = sql.NullTime{Valid: false}
			_, err := lastTask.Update()
			if err != nil {
				task.CompletedAt = sql.NullTime{Valid: true}
				c.infoPanel.Error(fmt.Errorf("failed to undo mark task as done: %s", err))
			} else {
				refresh()
			}
		// EDIT
		case 'e':
			c.editTaskForm(task, refresh, quit)
		// RESET TIMER
		case 'r':
			resetTimer()
			refresh()
		default:
			message := "invalid command: "
			if r == 0 || (string(r) != " " && unicode.IsSpace(r)) {
				c.infoPanel.Warn(message + event.Name())
			} else {
				c.infoPanel.Warn(message + "'" + string(r) + "'")
			}
		}
		return nil
	}
}

func (c *Controller) startDash(quit func()) {
	task, err := c.displayPanel.ShowTopTask()
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to start dash: %s", err))
		return
	}
	if lastTask.Id == 0 {
		lastTask = task
	}

	endAction := func() { trackTime(c, task, c.dashDuration) }
	timer := countdowntimer.Instance("lock in.", "dash complete. restart when ready.", c.dashDuration, endAction)

	timer.SetDescription(fmt.Sprintf(
		`([%[1]s::b]d[-:-:-]) mark task as done, ([%[1]s::b]u[-:-:-]) to undo
([%[1]s::b]e[-:-:-]) edit task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]b[-:-:-]) back`,
		tcell.ColorLimeGreen))

	timer.Layout.SetInputCapture(setDashCommands(c, quit, task, func() { timer.Reset(c.app) }))

	c.inputPanel.Set("Dash", timer.Layout)

	timer.Start(c.app)
}

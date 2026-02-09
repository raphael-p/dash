package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/database"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/countdowntimer"
)

var lastTask *database.Task

func trackTime(c *Controller, task *database.Task) func(time.Time, bool) {
	return func(lastRunTime time.Time, isEnd bool) {
		newTimeSpent := time.Duration(task.TimeSpentSeconds.Int16)*time.Second + time.Since(lastRunTime)
		task.TimeSpentSeconds = sql.NullInt16{
			Int16: int16(math.Round(newTimeSpent.Seconds())),
			Valid: true,
		}
		ok, err := task.Update()
		if !ok || err != nil {
			if err == nil {
				err = errors.New("noop")
			}
			c.infoPanel.Error(fmt.Errorf("failed to update time spent on task %d: %s", task.Id, err))
		}

		if isEnd && newTimeSpent > time.Minute {
			c.infoPanel.Info(fmt.Sprintf(
				"time tracking: %d minutes spent on task %d",
				int(newTimeSpent.Minutes()),
				task.Id,
			))
		}
	}
}

func setDashCommands(c *Controller, task *database.Task, resetTimer func(), quit func()) func(event *tcell.EventKey) *tcell.EventKey {
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
	t, err := c.displayPanel.ShowTopTask()
	task := &t
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to start dash: %s", err))
		return
	}
	if lastTask == nil {
		lastTask = task
	}

	timer := countdowntimer.Instance()
	timer.SetConfig(countdowntimer.Config{
		StartMessage:       "lock in.",
		EndMessage:         "dash complete. restart when ready.",
		CountdownDuration:  c.dashDuration,
		PeriodicSideEffect: trackTime(c, task),
	})

	timer.SetDescription(fmt.Sprintf(
		`([%[1]s::b]d[-:-:-]) mark task as done, ([%[1]s::b]u[-:-:-]) to undo
([%[1]s::b]e[-:-:-]) edit task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]b[-:-:-]) back`,
		tcell.ColorLimeGreen))

	timer.Layout.SetInputCapture(setDashCommands(c, task, func() { timer.Reset(c.app) }, quit))

	c.inputPanel.Set("Dash", timer.Layout)

	timer.Start(c.app)
}

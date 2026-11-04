package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/logger"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

var currentDashTask *database.Task
var lastDashTask *database.Task

func trackTime(c *Controller) func(time.Time, bool) {
	return func(lastRunTime time.Time, isEnd bool) {
		newTimeSpent := time.Duration(currentDashTask.TimeSpentSeconds.Int16)*time.Second + time.Since(lastRunTime)
		currentDashTask.TimeSpentSeconds = sql.NullInt16{
			Int16: int16(math.Round(newTimeSpent.Seconds())),
			Valid: true,
		}
		ok, err := currentDashTask.Update()
		if !ok || err != nil {
			if err == nil {
				err = errNoop
			}
			c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedDashUpdate, currentDashTask.ID, err))
		}

		if isEnd && newTimeSpent > time.Minute {
			c.infoPanel.Info(fmt.Sprintf(
				"time tracking: %d minutes spent on task [%d]",
				int(newTimeSpent.Minutes()),
				currentDashTask.ID,
			))
		}

		if isEnd {
			err = beeep.Alert("dash complete.", "good job.", "")
			if err != nil {
				logger.Warningf("failed to emit dash end alert: %v", err)
			}
		}
	}
}

func (c *Controller) openDash() {
	timer, isNew := tviewcomponents.CountdownTimer()

	t, err := database.GetTopTask()
	currentDashTask = &t
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedDashStart, err))
		}
		timer.Stop()
		c.Home()
		return
	}
	c.displayPanel.ShowTask(t)
	if lastDashTask == nil {
		lastDashTask = currentDashTask
	}

	menu := newInputKeybindMenu(c, config.KeybindInput, nil, func() { c.Home(); c.infoPanel.Clear() }).
		AddRuneKeybind('r', "restart timer", func() { timer.Reset(); timer.Start(c.queueUpdateDraw) }).
		AddRuneKeybind('d', "mark as done", markTaskDone(c))
	if lastDashTask.ID != currentDashTask.ID {
		menu.AddRuneKeybind('u', "unmark as done", unmarkTaskDone(c))
	}
	menu.AddRuneKeybind('a', "add task", func() { c.addTaskName(c.openDash) }).
		AddRuneKeybind('e', "edit task", func() { c.editTask(currentDashTask, c.openDash) })

	c.commandLine.Clear().Hide(menu.InputCapture())
	timer.SetText(menu.Text())
	c.navPanel.SetTitle("Dash").SetPrimitive(timer.Layout)

	if isNew {
		timer.SetConfig(tviewcomponents.TimerConfig{
			StartMessage:       "lock in.",
			EndMessage:         "dash complete. restart when ready.",
			CountdownDuration:  config.C.DashDuration,
			PeriodicSideEffect: trackTime(c),
		})
		timer.Start(c.queueUpdateDraw)
	}
}

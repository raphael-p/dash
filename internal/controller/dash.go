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
	timer, isNew := tviewcomponents.Timer()

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

	handler := &keybindHandler{c}
	mode := config.KeybindInput
	menu := tviewcomponents.NewKeybindMenu().
		AddKeybind('r', "restart timer", tviewcomponents.DefaultBind, func(_ rune) { timer.Reset(); timer.Start(c.queueUpdateDraw) }).
		AddKeybind('d', "mark task as done", tviewcomponents.DefaultBind, handler.markTaskDone).
		AddKeybind('u', "unmark task as done", tviewcomponents.DefaultBind, handler.unmarkTaskDone).
		AddKeybind('a', "add a new task", tviewcomponents.DefaultBind, handler.addTaskFlow(c.openDash)).
		AddKeybind('e', "edit task", tviewcomponents.DefaultBind, handler.editDashTask).
		AddKeybind('b', "back", tviewcomponents.DefaultBind, func(_ rune) { c.Home(); c.infoPanel.Clear() })

	charLimit, placeholder := config.GetCharLimit(mode), config.GetPlaceholder(mode)
	c.commandLine.Clear().Set(charLimit, placeholder, menu.InputCapture())
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

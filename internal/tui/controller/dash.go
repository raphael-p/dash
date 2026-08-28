package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/countdowntimer"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
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
				err = errors.New("noop")
			}
			c.infoPanel.Error(fmt.Errorf("failed to update time spent on task [%d]: %s", currentDashTask.ID, err))
		}

		if isEnd && newTimeSpent > time.Minute {
			c.infoPanel.Info(fmt.Sprintf(
				"time tracking: %d minutes spent on task [%d]",
				int(newTimeSpent.Minutes()),
				currentDashTask.ID,
			))
		}
	}
}

func (c *Controller) startDash() {
	t, err := c.displayPanel.ShowTopTask()
	currentDashTask = &t
	if err != nil {
		c.infoPanel.Error(fmt.Errorf("failed to start dash: %s", err))
		return
	}
	if lastDashTask == nil {
		lastDashTask = currentDashTask
	}

	timer := countdowntimer.Instance()
	timer.SetConfig(countdowntimer.Config{
		StartMessage:       "lock in.",
		EndMessage:         "dash complete. restart when ready.",
		CountdownDuration:  c.config.dashDuration,
		PeriodicSideEffect: trackTime(c),
	})

	handler := &keybindHandler{c}
	c.infoPanel.SetOnInputChange(func(_ string) {})
	keybindmenu.New().
		SetHighlighColour(tcell.ColorLimeGreen).
		SetFallback(handler.fallback).
		AddKeybind('d', "mark task as done", keybindmenu.DefaultBind, handler.markTaskDone).
		AddKeybind('u', "unmark task as done", keybindmenu.DefaultBind, handler.unmarkTaskDone).
		AddKeybind('e', "edit task", keybindmenu.DefaultBind, handler.editTask).
		AddKeybind('r', "restart timer", keybindmenu.DefaultBind, handler.resetTimer(func() { timer.Reset(c.app) })).
		AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.addFromDash).
		AddKeybind('b', "back", keybindmenu.DefaultBind, handler.backToHome).
		Apply(timer)

	c.inputPanel.Set("Dash", timer.Layout)
	timer.Start(c.app)
}

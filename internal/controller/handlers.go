package controller

import (
	"cmp"
	"database/sql"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/raphael-p/dash/internal/components"
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/sqlutil"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents"
)

func newInputKeybindMenu(c *Controller, mode config.InputMode, onEnter, onEsc tviewcomponents.KeybindHandler) *tviewcomponents.KeybindMenu {
	menu := tviewcomponents.NewKeybindMenu()
	switch mode {
	case config.DescriptionInput:
		menu.AddSecondaryKeybind("CTRL+J", "new line", tviewcomponents.BindCtrlJ, c.commandLine.NewLine)
		fallthrough
	case config.NameInput:
		menu.AddPassthroughKeybind(tviewcomponents.BindPrintable).
			AddSecondaryKeybind("CTRL+U", "clear line", tviewcomponents.BindCtrlU, c.commandLine.ClearLine)
	case config.TaskIDInput:
		menu.AddPassthroughKeybind(tviewcomponents.BindNumber)
	}
	if mode == config.NameInput || mode == config.DescriptionInput || mode == config.TaskIDInput {
		menu.AddPassthroughKeybind(tviewcomponents.BindBackspace)
	}

	if onEnter != nil {
		menu.AddPrimaryKeybind("ENTER", "submit", tviewcomponents.BindEnter, onEnter)
	}
	if onEsc == nil {
		reset := func() { c.displayPanel.ResetPagination(); c.infoPanel.Clear(); c.Home() }
		menu.AddSecondaryKeybind("ESC", "reset", tviewcomponents.BindEsc, reset).
			AddSecondaryKeybind("CTRL+C", "quit", tviewcomponents.BindRune, nil)
	} else {
		menu.AddSecondaryKeybind("ESC", "back", tviewcomponents.BindEsc, onEsc)
	}
	return menu
}

func askForConfirmation(c *Controller, action string, onYes, onNo func()) {
	menu := tviewcomponents.NewKeybindMenu().
		AddRuneKeybind('y', "yes", onYes).
		AddRuneKeybind('n', "no", onNo)
	c.commandLine.Clear().Hide(menu.InputCapture())
	c.navPanel.SetText(fmt.Sprintf("confirm action: %s?\n\n%s", action, menu.Text()))
}

func copy(c *Controller, s string) tviewcomponents.KeybindHandler {
	return func() {
		if err := clipboard.WriteAll(s); err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCopy, err))
		} else {
			c.infoPanel.Info("copied to clipboard")
		}
	}
}

func submitNameCandidate(c *Controller, back func()) tviewcomponents.KeybindHandler {
	return func() {
		if name := strutil.OrWarn(c.commandLine.GetInput(), "name", c.infoPanel.Warn); name != "" {
			c.addTaskDescription(name, back)
		}
	}
}

func submitDescriptionCandidate(c *Controller, name string, back func()) tviewcomponents.KeybindHandler {
	return func() {
		if description := strutil.OrWarn(c.commandLine.GetInput(), "name", c.infoPanel.Warn); description != "" {
			c.addTaskInsertionOrder(name, description, back)
		}
	}
}

func bumpTask(c *Controller, task *database.Task, onSuccess func()) tviewcomponents.KeybindHandler {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(c.commandLine.GetInput()))
		if task != nil {
			task.PriorityBumpedAt = sqlutil.NowNullTime()
			if service.UpdateTask(task, "priority bumped") {
				onSuccess()
			}
		}
	}
}

func submitName(c *Controller, task *database.Task, onSuccess func()) tviewcomponents.KeybindHandler {
	return func() {
		if name := strutil.OrWarn(c.commandLine.GetInput(), "name", c.infoPanel.Warn); name != "" {
			task.Name = name
			if database.Service().UpdateTask(task, "name changed") {
				onSuccess()
			}
		}
	}
}

func submitDescription(c *Controller, task *database.Task, onSuccess func()) tviewcomponents.KeybindHandler {
	return func() {
		if description := strutil.OrWarn(c.commandLine.GetInput(), "description", c.infoPanel.Warn); description != "" {
			task.Description = description
			if database.Service().UpdateTask(task, "description changed") {
				onSuccess()
			}
		}
	}
}

func openTask(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if id, ok := strutil.ParseIDOrWarn(c.commandLine.GetInput(), c.infoPanel.Warn); ok {
			if task, ok := database.Service().GetTask(id); ok {
				c.displayPanel.ShowTask(task)
				c.editTask(&task, c.Home)
			}
		}
	}
}

func reactivateTask(c *Controller, task *database.Task, onSuccess func()) tviewcomponents.KeybindHandler {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(c.commandLine.GetInput()))
		if task != nil {
			task.CompletedAt = sql.NullTime{Valid: false}
			if service.UpdateTask(task, "reactivated") {
				onSuccess()
			}
		}
	}
}

func removeTaskFromHome(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if id, ok := strutil.ParseIDOrWarn(c.commandLine.GetInput(), c.infoPanel.Warn); ok {
			askForConfirmation(
				c,
				fmt.Sprintf("delete task [%d]", id),
				func() { database.Service().RemoveTask(id); c.Home() },
				c.Home,
			)
		}
	}
}

func removeTaskFromEdit(c *Controller, task *database.Task, back func()) tviewcomponents.KeybindHandler {
	return func() {
		askForConfirmation(
			c,
			fmt.Sprintf("delete task [%d]", task.ID),
			func() { database.Service().RemoveTask(task.ID); back() },
			func() { c.editTask(task, back) },
		)
	}
}

func markTaskDone(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		currentDashTask.CompletedAt = sqlutil.NowNullTime()
		if database.Service().UpdateTask(currentDashTask, "done") {
			lastDashTask = currentDashTask
			c.openDash()
		} else {
			currentDashTask.CompletedAt = sql.NullTime{Valid: false}
		}
	}
}

func unmarkTaskDone(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if lastDashTask.ID == currentDashTask.ID {
			return // noop
		}
		lastDashTask.CompletedAt = sql.NullTime{Valid: false}

		if database.Service().UpdateTask(lastDashTask, "undone") {
			c.openDash()
		}
	}
}

func toggleTaskPause(c *Controller, task *database.Task, onSuccess func()) tviewcomponents.KeybindHandler {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(c.commandLine.GetInput()))
		if task != nil {
			if task.PausedAt.Valid {
				task.PausedAt = sql.NullTime{Valid: false}
				if service.UpdateTask(task, "unpaused") {
					onSuccess()
				}

			} else {
				task.PausedAt = sqlutil.NowNullTime()
				if service.UpdateTask(task, "paused") {
					onSuccess()
				}

			}
		}
	}
}

func scrollDownEdit(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if err := c.displayPanel.ScrollTask(components.ScrollDown); err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
		}
	}
}

func scrollUpEdit(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if err := c.displayPanel.ScrollTask(components.ScrollUp); err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
		}
	}
}

func scrollDownHome(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if err := c.displayPanel.ScrollTable(components.ScrollDown); err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
		}
	}
}

func scrollUpHome(c *Controller) tviewcomponents.KeybindHandler {
	return func() {
		if err := c.displayPanel.ScrollTable(components.ScrollUp); err != nil {
			c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
		}
	}
}

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

type keybindHandler struct {
	c *Controller
}

func (kh *keybindHandler) newMenu(mode config.InputMode, onEnter tviewcomponents.KeybindHandler) *tviewcomponents.KeybindMenu {
	menu := tviewcomponents.NewKeybindMenu()

	makeDescription := func(shortcut, function string) string {
		return fmt.Sprintf("[%s::b]%s[-:-:-]\t%s", config.C.HighlightColour, shortcut, function)
	}

	switch mode {
	case config.DescriptionInput:
		menu.AddKeybind(0, makeDescription("CTRL+J", "new line"), tviewcomponents.BindNewLine, kh.c.commandLine.NewLine)
		fallthrough
	case config.NameInput:
		menu.AddKeybind(0, "", tviewcomponents.BindPrintable, nil).
			AddKeybind(0, makeDescription("CTRL+U", "clear line"), tviewcomponents.BindClearLine, kh.c.commandLine.ClearLine)
	case config.TaskIDInput:
		menu.AddKeybind(0, "", tviewcomponents.BindNumber, nil)
	default:
		return menu
	}

	return menu.
		AddKeybind(0, makeDescription("ENTER", "submit"), tviewcomponents.BindEnter, onEnter).
		AddKeybind(0, "", tviewcomponents.BindBackspace, nil)
}

func (kh *keybindHandler) askForConfirmation(action string, onYes, onNo func()) {
	mode := config.KeybindInput
	menu := tviewcomponents.NewKeybindMenu().
		AddKeybind('y', "yes", tviewcomponents.DefaultBind, onYes).
		AddKeybind('n', "no", tviewcomponents.DefaultBind, onNo)
	kh.c.commandLine.Clear().Set(config.GetCharLimit(mode), config.GetPlaceholder(mode), nil)
	kh.c.navPanel.SetText(fmt.Sprintf("confirm action: %s?\n\n%s", action, menu.Text()))
}

func (kh *keybindHandler) copy(s string) func() {
	return func() {
		if err := clipboard.WriteAll(s); err != nil {
			kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCopy, err))
		} else {
			kh.c.infoPanel.Info("copied to clipboard")
		}
	}
}

func (kh *keybindHandler) addTaskFlowSubmitName(back func()) func() {
	return func() {
		if name := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); name != "" {
			kh.c.addTaskDescription(name, back)
		}
	}
}

func (kh *keybindHandler) addTaskFlowSubmitDescription(name string, back func()) func() {
	return func() {
		if description := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); description != "" {
			kh.c.addTaskInsertionOrder(name, description, back)
		}
	}
}

func (kh *keybindHandler) bumpTask(task *database.Task, onSuccess func()) func() {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(kh.c.commandLine.GetInput()))
		if task != nil {
			task.PriorityBumpedAt = sqlutil.NowNullTime()
			if service.UpdateTask(task, "priority bumped") {
				onSuccess()
			}
		}
	}
}

func (kh *keybindHandler) editTaskFlowSubmitName(task *database.Task, back func()) func() {
	return func() {
		if name := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); name != "" {
			service := database.Service()
			task.Name = name
			if service.UpdateTask(task, "name changed") {
				if task, ok := service.GetTask(task.ID); ok {
					kh.c.displayPanel.ShowTask(task)
					kh.c.editTaskFlow(&task, back)
				}
			}
		}
	}
}

func (kh *keybindHandler) editTaskFlowSubmitDescription(task *database.Task, back func()) func() {
	return func() {
		if description := strutil.OrWarn(kh.c.commandLine.GetInput(), "description", kh.c.infoPanel.Warn); description != "" {
			service := database.Service()
			task.Description = description
			if service.UpdateTask(task, "description changed") {
				if task, ok := service.GetTask(task.ID); ok {
					kh.c.displayPanel.ShowTask(task)
					kh.c.editTaskFlow(&task, back)
				}
			}
		}
	}
}

func (kh *keybindHandler) openTask() {
	if id, ok := strutil.ParseIDOrWarn(kh.c.commandLine.GetInput(), kh.c.infoPanel.Warn); ok {
		if task, ok := database.Service().GetTask(id); ok {
			kh.c.displayPanel.ShowTask(task)
			kh.c.editTaskFlow(&task, kh.c.Home)
		}
	}
}

func (kh *keybindHandler) reactivateTask(task *database.Task, onSuccess func()) func() {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(kh.c.commandLine.GetInput()))
		if task != nil {
			task.CompletedAt = sql.NullTime{Valid: false}
			if service.UpdateTask(task, "reactivated") {
				onSuccess()
			}
		}
	}
}

func (kh *keybindHandler) removeTaskFromHome() {
	if id, ok := strutil.ParseIDOrWarn(kh.c.commandLine.GetInput(), kh.c.infoPanel.Warn); ok {
		kh.askForConfirmation(
			fmt.Sprintf("delete task [%d]", id),
			func() { database.Service().RemoveTask(id); kh.c.Home() },
			kh.c.Home,
		)
	}
}

func (kh *keybindHandler) removeTaskFromEdit(task *database.Task, back func()) func() {
	return func() {
		kh.askForConfirmation(
			fmt.Sprintf("delete task [%d]", task.ID),
			func() { database.Service().RemoveTask(task.ID); back() },
			func() { kh.c.editTaskFlow(task, back) },
		)
	}
}

func (kh *keybindHandler) markTaskDone() {
	currentDashTask.CompletedAt = sqlutil.NowNullTime()
	if database.Service().UpdateTask(currentDashTask, "done") {
		lastDashTask = currentDashTask
		kh.c.openDash()
	} else {
		currentDashTask.CompletedAt = sql.NullTime{Valid: false}
	}
}

func (kh *keybindHandler) unmarkTaskDone() {
	if lastDashTask.ID == currentDashTask.ID {
		return // noop
	}
	lastDashTask.CompletedAt = sql.NullTime{Valid: false}

	if database.Service().UpdateTask(lastDashTask, "undone") {
		kh.c.openDash()
	}
}

func (kh *keybindHandler) toggleTaskPause(task *database.Task, onSuccess func()) func() {
	return func() {
		service := database.Service()
		task = cmp.Or(task, service.GetTaskFromString(kh.c.commandLine.GetInput()))
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

func (kh *keybindHandler) scrollDownTable() {
	err := kh.c.displayPanel.ScrollTable(components.ScrollDown)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (kh *keybindHandler) scrollUpTable() {
	err := kh.c.displayPanel.ScrollTable(components.ScrollUp)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

func (kh *keybindHandler) scrollDownTask() {
	err := kh.c.displayPanel.ScrollTask(components.ScrollDown)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (kh *keybindHandler) scrollUpTask() {
	err := kh.c.displayPanel.ScrollTask(components.ScrollUp)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

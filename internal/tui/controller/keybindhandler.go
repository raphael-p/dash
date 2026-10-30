package controller

import (
	"cmp"
	"database/sql"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/raphael-p/dash/internal/config"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/internal/tui/components"
	"github.com/raphael-p/dash/pkg/sqlutil"
	"github.com/raphael-p/dash/pkg/strutil"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

type keybindHandler struct {
	c *Controller
}

func (kh *keybindHandler) newMenu(mode config.InputMode, onEnter keybindmenu.Handler) *keybindmenu.KeybindMenu {
	menu := keybindmenu.New()

	makeDescription := func(shortcut, function string) string {
		return fmt.Sprintf("[%s::b]%s[-:-:-]\t%s", config.C.HighlightColour, shortcut, function)
	}

	switch mode {
	case config.DescriptionInput:
		menu.AddKeybind(0, makeDescription("CTRL+J", "new line"), keybindmenu.BindNewLine, kh.newLine)
		fallthrough
	case config.NameInput:
		menu.AddKeybind(0, "", keybindmenu.BindPrintable, nil).
			AddKeybind(0, makeDescription("CTRL+U", "clear line"), keybindmenu.BindClearLine, kh.clearLine)
	case config.TaskIDInput:
		menu.AddKeybind(0, "", keybindmenu.BindNumber, nil)
	default:
		return menu
	}

	return menu.
		AddKeybind(0, makeDescription("ENTER", "submit"), keybindmenu.BindEnter, onEnter).
		AddKeybind(0, "", keybindmenu.BindBackspace, nil)
}

func (kh *keybindHandler) backToHome(_ rune) {
	kh.c.Home()
}

func (kh *keybindHandler) askForConfirmation(action string, onYes, onNo func()) {
	mode := config.KeybindInput
	menu := keybindmenu.New()
	menu.AddKeybind('y', "yes", keybindmenu.DefaultBind, func(_ rune) { onYes() }).
		AddKeybind('n', "no", keybindmenu.DefaultBind, func(_ rune) { onNo() })
	kh.c.commandLine.Clear().Set(config.GetCharLimit(mode), config.GetPlaceholder(mode), nil)
	kh.c.navPanel.SetText(fmt.Sprintf("confirm action: %s?\n\n%s", action, menu.Text()))
}

func (kh *keybindHandler) copy(s string) func(rune) {
	return func(_ rune) {
		if err := clipboard.WriteAll(s); err != nil {
			kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCopy, err))
		} else {
			kh.c.infoPanel.Info("copied to clipboard")
		}
	}
}

func (kh *keybindHandler) clearLine(_ rune) {
	s := kh.c.commandLine.GetInput()
	i := strings.LastIndex(s, "\n")
	if i == -1 {
		kh.c.commandLine.SetInput("")
	} else {
		kh.c.commandLine.SetInput(s[:i])
	}
}

func (kh *keybindHandler) newLine(_ rune) {
	kh.c.commandLine.SetInput(kh.c.commandLine.GetInput() + "\n")
}

func (kh *keybindHandler) openDash(_ rune) {
	kh.c.openDash()
}

func (kh *keybindHandler) quit(_ rune) {
	kh.c.stop()
}

func (kh *keybindHandler) addTaskFlow(back func()) func(rune) {
	return func(_ rune) {
		kh.c.addTaskName(back)
	}
}

func (kh *keybindHandler) addTaskFlowSubmitName(back func()) func(rune) {
	return func(_ rune) {
		if name := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); name != "" {
			kh.c.addTaskDescription(name, back)
		}
	}
}

func (kh *keybindHandler) addTaskFlowSubmitDescription(name string, back func()) func(rune) {
	return func(_ rune) {
		if description := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); description != "" {
			kh.c.addTaskInsertionOrder(name, description, back)
		}
	}
}

func (kh *keybindHandler) addTaskFlowSubmit(name, description string, insertionOrder taskInsertionOrder, back func()) func(rune) {
	return func(_ rune) {
		if kh.c.addTask(name, description, insertionOrder) {
			back()
		}
	}
}

func (kh *keybindHandler) bumpTask(task *database.Task, onSuccess func()) func(rune) {
	return func(_ rune) {
		task = cmp.Or(task, kh.c.getTaskFromCommandLine())
		if task != nil {
			task.PriorityBumpedAt = sqlutil.NowNullTime()
			if kh.c.updateTask(task, "priority bumped") {
				onSuccess()
			}
		}
	}
}

func (kh *keybindHandler) editDashTask(_ rune) {
	kh.c.editTaskFlow(currentDashTask, kh.c.openDash)
}

func (kh *keybindHandler) editTaskFlowName(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		kh.c.editTaskName(task, back)
	}
}

func (kh *keybindHandler) editTaskFlowDescription(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		kh.c.editTaskDescription(task, back)
	}
}

func (kh *keybindHandler) editTaskFlowSubmitName(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		if name := strutil.OrWarn(kh.c.commandLine.GetInput(), "name", kh.c.infoPanel.Warn); name != "" {
			task.Name = name
			if kh.c.updateTask(task, "name changed") {
				kh.c.openTask(task.ID, back)
			}
		}
	}
}

func (kh *keybindHandler) editTaskFlowSubmitDescription(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		if description := strutil.OrWarn(kh.c.commandLine.GetInput(), "description", kh.c.infoPanel.Warn); description != "" {
			task.Description = description
			if kh.c.updateTask(task, "description changed") {
				kh.c.openTask(task.ID, back)
			}
		}
	}
}

func (kh *keybindHandler) openTask(_ rune) {
	if id, ok := strutil.ParseIDOrWarn(kh.c.commandLine.GetInput(), kh.c.infoPanel.Warn); ok {
		kh.c.openTask(id, kh.c.Home)
	}
}

func (kh *keybindHandler) reactivateTask(task *database.Task, onSuccess func()) func(rune) {
	return func(_ rune) {
		task = cmp.Or(task, kh.c.getTaskFromCommandLine())
		if task != nil {
			task.CompletedAt = sql.NullTime{Valid: false}
			if kh.c.updateTask(task, "reactivated") {
				onSuccess()
			}
		}
	}
}

func (kh *keybindHandler) removeTaskFromHome(_ rune) {
	if id, ok := strutil.ParseIDOrWarn(kh.c.commandLine.GetInput(), kh.c.infoPanel.Warn); ok {
		kh.askForConfirmation(
			fmt.Sprintf("delete task [%d]", id),
			func() { kh.c.removeTask(id); kh.c.Home() },
			kh.c.Home,
		)
	}
}

func (kh *keybindHandler) removeTaskFromEdit(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		kh.askForConfirmation(
			fmt.Sprintf("delete task [%d]", task.ID),
			func() { kh.c.removeTask(task.ID); back() },
			func() { kh.c.editTaskFlow(task, back) },
		)
	}
}

func (kh *keybindHandler) markTaskDone(_ rune) {
	currentDashTask.CompletedAt = sqlutil.NowNullTime()
	if kh.c.updateTask(currentDashTask, "done") {
		lastDashTask = currentDashTask
		kh.c.openDash()
	} else {
		currentDashTask.CompletedAt = sql.NullTime{Valid: false}
	}
}

func (kh *keybindHandler) unmarkTaskDone(_ rune) {
	if lastDashTask.ID == currentDashTask.ID {
		return // noop
	}
	lastDashTask.CompletedAt = sql.NullTime{Valid: false}

	if kh.c.updateTask(lastDashTask, "undone") {
		kh.c.openDash()
	}
}

func (kh *keybindHandler) toggleTaskPause(task *database.Task, onSuccess func()) func(rune) {
	return func(_ rune) {
		task = cmp.Or(task, kh.c.getTaskFromCommandLine())
		if task != nil {
			if task.PausedAt.Valid {
				task.PausedAt = sql.NullTime{Valid: false}
				if kh.c.updateTask(task, "unpaused") {
					onSuccess()
				}

			} else {
				task.PausedAt = sqlutil.NowNullTime()
				if kh.c.updateTask(task, "paused") {
					onSuccess()
				}

			}
		}
	}
}

func (kh *keybindHandler) toggleTaskMode(_ rune) {
	kh.c.displayPanel.ToggleTaskMode()
	kh.c.Home()
}

func (kh *keybindHandler) scrollDownTable(_ rune) {
	err := kh.c.displayPanel.ScrollTable(components.ScrollDown)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (kh *keybindHandler) scrollUpTable(_ rune) {
	err := kh.c.displayPanel.ScrollTable(components.ScrollUp)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

func (kh *keybindHandler) scrollDownTask(_ rune) {
	err := kh.c.displayPanel.ScrollTask(components.ScrollDown)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (kh *keybindHandler) scrollUpTask(_ rune) {
	err := kh.c.displayPanel.ScrollTask(components.ScrollUp)
	if err != nil {
		kh.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

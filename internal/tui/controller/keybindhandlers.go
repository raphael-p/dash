package controller

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/internal/tui/panels"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

type keybindHandler struct {
	c *Controller
}

func (ka *keybindHandler) backToHome(_ rune) {
	ka.c.Home()
}

func (ka *keybindHandler) askForConfirmation(action string, onYes, onNo func()) {
	handler := &keybindHandler{ka.c}
	menu := configureKeybindMenu(handler, nil, keybindInput)
	menu.AddKeybind('y', "yes", keybindmenu.DefaultBind, func(_ rune) { onYes() }).
		AddKeybind('n', "no", keybindmenu.DefaultBind, func(_ rune) { onNo() })
	configureCommandLine(ka.c, menu, keybindInput)
	ka.c.navPanel.SetText(fmt.Sprintf("confirm action: %s?\n\n%s", action, menu.GenerateText()))
}

func (ka *keybindHandler) copy(s string) func(rune) {
	return func(_ rune) {
		if err := clipboard.WriteAll(s); err != nil {
			ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedCopy, err))
		} else {
			ka.c.infoPanel.Info("copied to clipboard")
		}
	}
}

func (ka *keybindHandler) clearLine(_ rune) {
	s := ka.c.commandLine.GetInput()
	i := strings.LastIndex(s, "\n")
	if i == -1 {
		ka.c.commandLine.SetInput("")
	} else {
		ka.c.commandLine.SetInput(s[:i])
	}
}

func (ka *keybindHandler) fallback(commandName string) {
	ka.c.infoPanel.Warn("invalid command: " + commandName)
}

func (ka *keybindHandler) newLine(_ rune) {
	ka.c.commandLine.SetInput(ka.c.commandLine.GetInput() + "\n")
}

func (ka *keybindHandler) openDash(_ rune) {
	ka.c.openDash()
}

func (ka *keybindHandler) quit(_ rune) {
	ka.c.app.Stop()
}

func (ka *keybindHandler) addTaskFlow(back func()) func(rune) {
	return func(_ rune) {
		ka.c.addTaskName(back)
	}
}

func (ka *keybindHandler) addTaskFlowSubmitName(back func()) func(rune) {
	return func(_ rune) {
		if name := getInputOrWarn(ka.c.commandLine.GetInput(), "name", ka.c.infoPanel.Warn); name != "" {
			ka.c.addTaskDescription(name, back)
		}
	}
}

func (ka *keybindHandler) addTaskFlowSubmitDescription(name string, back func()) func(rune) {
	return func(_ rune) {
		if description := getInputOrWarn(ka.c.commandLine.GetInput(), "name", ka.c.infoPanel.Warn); description != "" {
			ka.c.addTaskInsertionOrder(name, description, back)
		}
	}
}

func (ka *keybindHandler) addTaskFlowSubmit(name, description string, insertionOrder taskInsertionOrder, back func()) func(rune) {
	return func(_ rune) {
		if ka.c.addTask(name, description, insertionOrder) {
			back()
		}
	}
}

func (ka *keybindHandler) bumpTaskFromHome(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		if ka.c.bumpTask(id) {
			ka.c.refreshTasks()
			ka.c.commandLine.Clear()
		}
	}
}

func (ka *keybindHandler) bumpTaskFromEdit(id int64) func(rune) {
	return func(_ rune) {
		ka.c.bumpTask(id)
	}
}

func (ka *keybindHandler) editDashTask(_ rune) {
	ka.c.editTaskFlow(currentDashTask, ka.c.openDash)
}

func (ka *keybindHandler) editTaskFlowName(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		ka.c.editTaskName(task, back)
	}
}

func (ka *keybindHandler) editTaskFlowDescription(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		ka.c.editTaskDescription(task, back)
	}
}

func (ka *keybindHandler) editTaskFlowSubmitName(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		if name := getInputOrWarn(ka.c.commandLine.GetInput(), "name", ka.c.infoPanel.Warn); name != "" {
			if ka.c.editTask(task, name, "") {
				ka.c.openTask(task.ID, back)
			}
		}
	}
}

func (ka *keybindHandler) editTaskFlowSubmitDescription(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		if description := getInputOrWarn(ka.c.commandLine.GetInput(), "description", ka.c.infoPanel.Warn); description != "" {
			if ka.c.editTask(task, "", description) {
				ka.c.openTask(task.ID, back)
			}
		}
	}
}

func (ka *keybindHandler) openTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		ka.c.openTask(id, ka.c.Home)
	}
}

func (ka *keybindHandler) reactivateTaskFromHome(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		if ka.c.reactivateTask(id) {
			ka.c.refreshTasks()
			ka.c.commandLine.Clear()
		}
	}
}

func (ka *keybindHandler) reactivateTaskFromEdit(id int64) func(rune) {
	return func(_ rune) {
		ka.c.reactivateTask(id)
	}
}

func (ka *keybindHandler) removeTaskFromHome(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		ka.askForConfirmation(
			fmt.Sprintf("delete task [%d]", id),
			func() { ka.c.removeTask(id); ka.c.Home() },
			ka.c.Home,
		)
	}
}

func (ka *keybindHandler) removeTaskFromEdit(task *database.Task, back func()) func(rune) {
	return func(_ rune) {
		ka.askForConfirmation(
			fmt.Sprintf("delete task [%d]", task.ID),
			func() { ka.c.removeTask(task.ID); back() },
			func() { ka.c.editTaskFlow(task, back) },
		)
	}
}

func (ka *keybindHandler) markTaskDone(_ rune) {
	currentDashTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	_, err := currentDashTask.Update()
	if err != nil {
		currentDashTask.CompletedAt = sql.NullTime{Valid: false}
		ka.c.infoPanel.Error(fmt.Errorf("%w [%d]: %w", errFailedMark, currentDashTask.ID, err))
	} else {
		lastDashTask = currentDashTask
		ka.c.openDash()
	}
}

func (ka *keybindHandler) unmarkTaskDone(_ rune) {
	if lastDashTask.ID == currentDashTask.ID {
		return // noop
	}
	lastDashTask.CompletedAt = sql.NullTime{Valid: false}
	_, err := lastDashTask.Update()
	if err != nil {
		currentDashTask.CompletedAt = sql.NullTime{Valid: true}
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedUnmark, err))
	} else {
		ka.c.openDash()
	}
}

func (ka *keybindHandler) toggleTaskMode(_ rune) {
	ka.c.displayPanel.ToggleTaskMode()
	ka.c.Home()
}

func (ka *keybindHandler) scrollDownTable(_ rune) {
	err := ka.c.displayPanel.ScrollTable(panels.ScrollDown)
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (ka *keybindHandler) scrollUpTable(_ rune) {
	err := ka.c.displayPanel.ScrollTable(panels.ScrollUp)
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

func (ka *keybindHandler) scrollDownTask(_ rune) {
	err := ka.c.displayPanel.ScrollTask(panels.ScrollDown)
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (ka *keybindHandler) scrollUpTask(_ rune) {
	err := ka.c.displayPanel.ScrollTask(panels.ScrollUp)
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

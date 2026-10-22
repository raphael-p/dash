package controller

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/raphael-p/dash/internal/database"
)

type keybindHandler struct {
	c *Controller
}

func (ka *keybindHandler) backToHome(_ rune) {
	ka.c.Home()
}

func (ka *keybindHandler) openDash(_ rune) {
	ka.c.openDash()
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

func (ka *keybindHandler) copyTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		ka.c.copyTask(id)
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
		if ka.c.removeTask(id) {
			ka.c.refreshTasks()
			ka.c.commandLine.Clear()
		}
	}
}

func (ka *keybindHandler) removeTaskFromEdit(id int64, back func()) func(rune) {
	return func(_ rune) {
		if ka.c.removeTask(id) {
			back()
		}
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

func (ka *keybindHandler) fallback(commandName string) {
	ka.c.infoPanel.Warn("invalid command: " + commandName)
}

func (ka *keybindHandler) quit(_ rune) {
	ka.c.app.Stop()
}

func (ka *keybindHandler) scrollDown(_ rune) {
	err := ka.c.displayPanel.ScrollDown()
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollDown, err))
	}
}

func (ka *keybindHandler) scrollUp(_ rune) {
	err := ka.c.displayPanel.ScrollUp()
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("%w: %w", errFailedScrollUp, err))
	}
}

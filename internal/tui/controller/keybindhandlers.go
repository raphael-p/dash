package controller

import (
	"database/sql"
	"fmt"
	"time"
)

type keybindHandler struct {
	c *Controller
}

func (ka *keybindHandler) addFromDash(_ rune) {
	ka.c.addTaskName(ka.c.startDash)
}

func (ka *keybindHandler) addFromHome(_ rune) {
	ka.c.addTaskName(ka.c.Home)
}

func (ka *keybindHandler) backToHome(_ rune) {
	ka.c.Home()
}

func (ka *keybindHandler) startDash(_ rune) {
	ka.c.startDash()
}

func (ka *keybindHandler) bumpTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		if ka.c.bumpTask(id) {
			ka.c.refreshTasks()
		}
	}
}

func (ka *keybindHandler) copyTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		ka.c.copyTask(id)
	}
}

func (ka *keybindHandler) editTask(_ rune) {
	ka.c.editTaskForm(currentDashTask, ka.c.startDash)
}

func (ka *keybindHandler) openTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		ka.c.openTask(id, ka.c.Home)
	}
}

func (ka *keybindHandler) reactivateTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		if ka.c.reactivateTask(id) {
			ka.c.refreshTasks()
		}
	}

}

func (ka *keybindHandler) removeTask(_ rune) {
	if id, ok := getIDOrWarn(ka.c.commandLine.GetInput(), ka.c.infoPanel.Warn); ok {
		if ka.c.removeTask(id) {
			ka.c.refreshTasks()
		}
	}
}

func (ka *keybindHandler) markTaskDone(_ rune) {
	currentDashTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	_, err := currentDashTask.Update()
	if err != nil {
		currentDashTask.CompletedAt = sql.NullTime{Valid: false}
		ka.c.infoPanel.Error(fmt.Errorf("failed to mark task [%d] as done: %s", currentDashTask.ID, err))
	} else {
		lastDashTask = currentDashTask
		ka.c.startDash()
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
		ka.c.infoPanel.Error(fmt.Errorf("failed to undo mark task as done: %s", err))
	} else {
		ka.c.startDash()
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

func (ka *keybindHandler) resetTimer(resetTimer func()) func(rune) {
	return func(_ rune) {
		resetTimer()
		ka.c.startDash()
	}
}

func (ka *keybindHandler) scrollDown(_ rune) {
	err := ka.c.displayPanel.ScrollDown()
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("failed to scroll down: %s", err))
	}
}

func (ka *keybindHandler) scrollUp(_ rune) {
	err := ka.c.displayPanel.ScrollUp()
	if err != nil {
		ka.c.infoPanel.Error(fmt.Errorf("failed to scroll up: %s", err))
	}
}

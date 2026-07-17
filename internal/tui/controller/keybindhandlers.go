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
	ka.c.addTaskForm(ka.c.startDash)
}

func (ka *keybindHandler) addFromHome(_ rune) {
	ka.c.addTaskForm(ka.c.Home)
}

func (ka *keybindHandler) backToHome(_ rune) {
	ka.c.Home()
}

func (ka *keybindHandler) backspace(r rune) {
	currentInput := ka.c.infoPanel.GetInput()
	if currentInput != "" {
		ka.c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
	}
}

func (ka *keybindHandler) bumpTask(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.bumpTask(ka.c.infoPanel.GetInput())
	} else {
		ka.c.bumpTaskForm(ka.c.Home)
	}
}

func (ka *keybindHandler) copyTask(_ rune) {
	ka.c.copyTask(ka.c.infoPanel.GetInput())
}

func (ka *keybindHandler) editTask(_ rune) {
	ka.c.editTaskForm(currentDashTask, ka.c.startDash)
}

func (ka *keybindHandler) startDash(_ rune) {
	ka.c.startDash()
}

func (ka *keybindHandler) markTaskDone(_ rune) {
	currentDashTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	_, err := currentDashTask.Update()
	if err != nil {
		currentDashTask.CompletedAt = sql.NullTime{Valid: false}
		ka.c.infoPanel.Error(fmt.Errorf("failed to mark task [%d] as done: %s", currentDashTask.Id, err))
	} else {
		lastDashTask = currentDashTask
		ka.c.startDash()
	}
}

func (ka *keybindHandler) fallback(commandName string) {
	ka.c.infoPanel.Warn("invalid command: " + commandName)
}

func (ka *keybindHandler) numberInput(r rune) {
	ka.c.infoPanel.SetInput(ka.c.infoPanel.GetInput() + string(r))
}

func (ka *keybindHandler) openTask(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.openTask(ka.c.infoPanel.GetInput(), ka.c.Home)
	} else {
		ka.c.openTaskForm(ka.c.Home)
	}
}

func (ka *keybindHandler) quit(_ rune) {
	ka.c.app.Stop()
}

func (ka *keybindHandler) removeTask(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.removeTask(ka.c.infoPanel.GetInput())
	} else {
		ka.c.removeTaskForm(ka.c.Home)
	}
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

func (ka *keybindHandler) unmarkTaskDone(_ rune) {
	if lastDashTask.Id == currentDashTask.Id {
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

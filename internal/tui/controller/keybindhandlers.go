package controller

import (
	"fmt"
)

type keybindHandler struct {
	c       *Controller
	refresh func()
}

func (ka *keybindHandler) add(_ rune) {
	ka.c.addTaskForm(ka.c.Home, ka.c.app.Stop)
}

func (ka *keybindHandler) backspace(r rune) {
	currentInput := ka.c.infoPanel.GetInput()
	if currentInput != "" {
		ka.c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
	}
	ka.refresh()
}

func (ka *keybindHandler) bump(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.bumpTask(ka.c.infoPanel.GetInput())
	} else {
		ka.c.bumpTaskForm(ka.c.Home, ka.c.app.Stop)
	}
}

func (ka *keybindHandler) dash(_ rune) {
	ka.c.startDash(ka.c.app.Stop)
}

func (ka *keybindHandler) numberInput(r rune) {
	currentInput := ka.c.infoPanel.GetInput()
	ka.c.infoPanel.SetInput(currentInput + string(r))
	ka.refresh()
}

func (ka *keybindHandler) open(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.openTask(ka.c.infoPanel.GetInput(), ka.c.Home, ka.c.app.Stop)
	} else {
		ka.c.openTaskForm(ka.c.Home, ka.c.app.Stop)
	}
}

func (ka *keybindHandler) quit(_ rune) {
	ka.c.app.Stop()
}

func (ka *keybindHandler) remove(_ rune) {
	if ka.c.infoPanel.GetInput() != "" {
		ka.c.removeTask(ka.c.infoPanel.GetInput())
	} else {
		ka.c.removeTaskForm(ka.c.Home, ka.c.app.Stop)
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

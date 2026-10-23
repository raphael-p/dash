package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/database"
	"github.com/raphael-p/dash/pkg/tviewcomponents/keybindmenu"
)

type commandLineMode int

const (
	descriptionInput commandLineMode = iota
	keybindInput
	nameInput
	taskIDInput
)

func configureCommandLine(c *Controller, menu *keybindmenu.KeybindMenu, mode commandLineMode) {
	var charLimit int
	var placeholder string
	switch mode {
	case descriptionInput:
		charLimit = c.config.descriptionCharLimit
		placeholder = "enter a task description, press ENTER to submit, press SHIFT+ENTER or CTRL+J to create a new line"
	case keybindInput:
		charLimit = 1
	case nameInput:
		charLimit = c.config.nameCharLimit
		placeholder = "enter a task name, press ENTER to submit"
	case taskIDInput:
		charLimit = 9
		placeholder = "enter a task id, press ENTER to submit"
	}
	c.commandLine.Reset().SetCharLimit(charLimit).SetPlaceholder(placeholder)
	if menu != nil {
		c.commandLine.SetInputCapture(menu.GenerateInputCapture())
	}
}

func configureKeybindMenu(handler *keybindHandler, onEnter func(rune), mode commandLineMode) *keybindmenu.KeybindMenu {
	menu := keybindmenu.New().SetHighlightColour(tcell.ColorLimeGreen).SetFallback(handler.fallback)

	switch mode {
	case descriptionInput:
		menu.AddKeybind(0, "", keybindmenu.BindNewLine, handler.newLine)
		fallthrough
	case nameInput:
		menu.AddKeybind(0, "", keybindmenu.BindPrintable, nil).
			AddKeybind(0, "", keybindmenu.BindClearLine, handler.clearLine)
	case taskIDInput:
		menu.AddKeybind(0, "", keybindmenu.BindNumber, nil)
	case keybindInput:
		return menu
	}

	return menu.
		AddKeybind(0, "", keybindmenu.BindEnter, onEnter).
		AddKeybind(0, "", keybindmenu.BindBackspace, nil)
}

func getIDOrWarn(idString string, warn func(err string)) (int64, bool) {
	id, err := extractIDFromString(idString)
	if err != nil {
		warn(fmt.Sprint("your input is invalid: ", err))
		return 0, false
	}
	return id, true
}

func extractIDFromString(idString string) (int64, error) {
	if idString == "" {
		return 0, errEmptyID
	}

	id, err := strconv.ParseInt(idString, 10, 0)
	if err != nil || id < 1 {
		return 0, errIDNotInt
	}

	if id < 1 {
		return 0, database.ErrIDNotPositive
	}
	return id, nil
}

func getInputOrWarn(input, label string, warn func(err string)) string {
	if value := strings.TrimSpace(input); value != "" {
		return value
	}
	warn(fmt.Sprintf("please provide a %s", label))
	return ""
}

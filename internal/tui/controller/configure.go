package controller

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/config"
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
		charLimit = config.C.DescriptionCharLimit
		placeholder = "enter a task description"
	case keybindInput:
		charLimit = 1
	case nameInput:
		charLimit = config.C.NameCharLimit
		placeholder = "enter a task name"
	case taskIDInput:
		charLimit = 9
		placeholder = "enter a task id"
	}
	c.commandLine.Reset().SetCharLimit(charLimit).SetPlaceholder(placeholder)
	if menu != nil {
		c.commandLine.SetInputCapture(menu.GenerateInputCapture())
	}
}

func configureKeybindMenu(handler *keybindHandler, onEnter func(rune), mode commandLineMode) *keybindmenu.KeybindMenu {
	menu := keybindmenu.New().SetHighlightColour(tcell.ColorLimeGreen).SetFallback(handler.fallback)

	makeDescription := func(shortcut, function string) string {
		return fmt.Sprintf("[%s::b]%s[-:-:-]\t%s", tcell.ColorLimeGreen, shortcut, function)
	}
	switch mode {
	case descriptionInput:
		menu.AddKeybind(0, makeDescription("CTRL+J", "new line"), keybindmenu.BindNewLine, handler.newLine)
		fallthrough
	case nameInput:
		menu.AddKeybind(0, "", keybindmenu.BindPrintable, nil).
			AddKeybind(0, makeDescription("CTRL+U", "clear line"), keybindmenu.BindClearLine, handler.clearLine)
	case taskIDInput:
		menu.AddKeybind(0, "", keybindmenu.BindNumber, nil)
	case keybindInput:
		return menu
	}

	return menu.
		AddKeybind(0, makeDescription("ENTER", "submit"), keybindmenu.BindEnter, onEnter).
		AddKeybind(0, "", keybindmenu.BindBackspace, nil)
}

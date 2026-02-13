package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.refreshTasks()
	handler := keybindHandler{c}

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)

	keybindmenu.New().SetHighlighColour(tcell.ColorLimeGreen).
		SetFallback(func(commandName string) { c.infoPanel.Warn("invalid command: " + commandName) }).
		AddKeybind('d', "dash", keybindmenu.DefaultBind, handler.dash).
		AddKeybind('o', "open a task", keybindmenu.BindEnter, handler.open).
		AddKeybind('b', "bump the priority of a task", keybindmenu.DefaultBind, handler.bump).
		AddKeybind('r', "remove a task", keybindmenu.DefaultBind, handler.remove).
		AddKeybind('a', "add a new task", keybindmenu.DefaultBind, handler.add).
		AddKeybind('j', "scroll down", keybindmenu.DefaultBind, handler.scrollDown).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, handler.scrollUp).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, handler.quit).
		AddKeybind(0, "or enter a task ID...", keybindmenu.BindNumber, handler.numberInput).
		AddKeybind(0, "", keybindmenu.BindBackspace, handler.backspace).
		SetFootnote("tip: you can type the task ID before invoking the open (o), bump (b), or remove (r) commands").
		Apply(home)

	c.inputPanel.Set("Welcome to Dash!", home)
}

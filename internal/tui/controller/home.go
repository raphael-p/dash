package controller

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	keybindmenu "github.com/raphael-p/datashard/pkg/tviewcomponents/keybindMenu"
	"github.com/rivo/tview"
)

var emptyTrigger = func(e tcell.Key, r rune) bool { return false }
var triggerOnEnter = func(e tcell.Key, r rune) bool {
	return e == tcell.KeyEnter
}
var triggerOnBackspace = func(e tcell.Key, r rune) bool {
	return e == tcell.KeyBackspace || e == tcell.KeyBackspace2
}
var triggerOnNumber = func(e tcell.Key, r rune) bool {
	return r >= '0' && r <= '9'
}

func (c *Controller) Home() {
	c.refreshTasks()
	back := c.Home
	quit := c.app.Stop

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)

	keybindmenu.New().SetHighlighColour(tcell.ColorLimeGreen).
		SetDefaultAction(func(name string, r rune) {
			message := "invalid command: "
			if r == 0 || (string(r) != " " && unicode.IsSpace(r)) {
				c.infoPanel.Warn(message + name)
			} else {
				c.infoPanel.Warn(message + "'" + string(r) + "'")
			}

		}).
		AddKeybind('d', "dash", emptyTrigger, func(r rune) { c.addTaskForm(back, quit) }).
		AddKeybind('o', "open a task", triggerOnEnter,
			func(r rune) {
				if c.infoPanel.GetInput() != "" {
					c.openTask(c.infoPanel.GetInput(), back, quit)
				} else {
					c.openTaskForm(back, quit)
				}
			}).
		AddKeybind('b', "bump the priority of a task", emptyTrigger,
			func(r rune) {
				if c.infoPanel.GetInput() != "" {
					c.bumpTask(c.infoPanel.GetInput())
				} else {
					c.bumpTaskForm(back, quit)
				}
			}).
		AddKeybind('r', "remove a task", emptyTrigger,
			func(r rune) {
				if c.infoPanel.GetInput() != "" {
					c.removeTask(c.infoPanel.GetInput())
				} else {
					c.removeTaskForm(back, quit)
				}
			}).
		AddKeybind('a', "add a new task", emptyTrigger, func(r rune) { c.addTaskForm(back, quit) }).
		AddKeybind('j', "scroll down", emptyTrigger,
			func(r rune) {
				err := c.displayPanel.ScrollDown()
				if err != nil {
					c.infoPanel.Error(fmt.Errorf("failed to scroll down: %s", err))
				}
			}).
		AddKeybind('k', "scroll up", emptyTrigger,
			func(r rune) {
				err := c.displayPanel.ScrollUp()
				if err != nil {
					c.infoPanel.Error(fmt.Errorf("failed to scroll up: %s", err))
				}
			}).
		AddKeybind('q', "quit", emptyTrigger, func(r rune) { quit() }).
		AddKeybind(0, "", triggerOnNumber, func(r rune) { c.infoPanel.SetInput(c.infoPanel.GetInput() + string(r)) }).
		AddKeybind(0, "", triggerOnBackspace,
			func(r rune) {
				currentInput := c.infoPanel.GetInput()
				if currentInput != "" {
					c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
				}
			}).
		SetFootnote("tip: you can type the task ID before invoking the open (o), bump (b), or remove (r) commands").
		Apply(home)

	c.inputPanel.Set("Welcome to Dash!", home)
}

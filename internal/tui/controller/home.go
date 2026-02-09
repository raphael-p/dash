package controller

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/pkg/tviewcomponents/keybindmenu"
	"github.com/rivo/tview"
)

func (c *Controller) Home() {
	c.refreshTasks()
	back := c.Home
	quit := c.app.Stop

	home := tview.NewTextView().SetDynamicColors(true)
	home.SetBorderPadding(1, 1, 2, 2)

	keybindmenu.New().SetHighlighColour(tcell.ColorLimeGreen).
		SetDefaultAction(func(commandName string) { c.infoPanel.Warn("invalid command: " + commandName) }).
		AddKeybind('d', "dash", keybindmenu.DefaultBind, func(_ rune) {
			c.startDash(c.app.Stop)
		}).
		AddKeybind('o', "open a task", keybindmenu.BindEnter, func(_ rune) {
			if c.infoPanel.GetInput() != "" {
				c.openTask(c.infoPanel.GetInput(), c.Home, c.app.Stop)
			} else {
				c.openTaskForm(c.Home, c.app.Stop)
			}
		}).
		AddKeybind('b', "bump the priority of a task", keybindmenu.DefaultBind, func(_ rune) {
			if c.infoPanel.GetInput() != "" {
				c.bumpTask(c.infoPanel.GetInput())
			} else {
				c.bumpTaskForm(back, quit)
			}
		}).
		AddKeybind('r', "remove a task", keybindmenu.DefaultBind, func(_ rune) {
			if c.infoPanel.GetInput() != "" {
				c.removeTask(c.infoPanel.GetInput())
			} else {
				c.removeTaskForm(back, quit)
			}
		}).
		AddKeybind('a', "add a new task", keybindmenu.DefaultBind, func(_ rune) {
			c.addTaskForm(back, quit)
		}).
		AddKeybind('j', "scroll down", keybindmenu.DefaultBind, func(_ rune) {
			err := c.displayPanel.ScrollDown()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to scroll down: %s", err))
			}
		}).
		AddKeybind('k', "scroll up", keybindmenu.DefaultBind, func(_ rune) {
			err := c.displayPanel.ScrollUp()
			if err != nil {
				c.infoPanel.Error(fmt.Errorf("failed to scroll up: %s", err))
			}
		}).
		AddKeybind('q', "quit", keybindmenu.DefaultBind, func(_ rune) {
			quit()
		}).
		AddKeybind(0, "", keybindmenu.BindNumber, func(r rune) {
			c.infoPanel.SetInput(c.infoPanel.GetInput() + string(r))
		}).
		AddKeybind(0, "", keybindmenu.BindBackspace, func(r rune) {
			currentInput := c.infoPanel.GetInput()
			if currentInput != "" {
				c.infoPanel.SetInput(currentInput[:len(currentInput)-1])
			}
		}).
		SetFootnote("tip: you can type the task ID before invoking the open (o), bump (b), or remove (r) commands").
		Apply(home)

	c.inputPanel.Set("Welcome to Dash!", home)
}

package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/datashard/internal/tui/controller"
	"github.com/raphael-p/datashard/internal/tui/panels"
	"github.com/raphael-p/datashard/pkg/logger"
	"github.com/rivo/tview"
)

func startTUI() {
	app := tview.NewApplication()
	displayPanel := panels.NewDisplayPanel()
	infoPanel := panels.NewInfoPanel()
	inputPanel := panels.NewInputPanel()

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputPanel.GetPanel(), 0, 3, true).
		AddItem(infoPanel.GetPanel(), 0, 1, false)
	flex := tview.NewFlex().
		AddItem(displayPanel.GetPanel(), 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	c := controller.NewController(app, inputPanel, infoPanel, displayPanel, time.Duration(config.DashDuration)*time.Second)
	c.Home()

	inputPanel.GetPanel().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			displayPanel.ResetPagination()
			c.Home()
			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/tui/controller"
	"github.com/raphael-p/dash/internal/tui/panels"
	"github.com/raphael-p/dash/pkg/logger"
	"github.com/rivo/tview"
)

func startTUI() {
	app := tview.NewApplication()
	displayPanel := panels.NewDisplayPanel()
	infoPanel := panels.NewInfoPanel()
	inputPanel := panels.NewInputPanel()

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputPanel.GetPrimitive(), 0, 3, true).
		AddItem(infoPanel.GetPrimitive(), 0, 1, false)
	flex := tview.NewFlex().
		AddItem(displayPanel.GetPrimitive(), 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	c := controller.NewController(
		app, inputPanel, infoPanel, displayPanel,
		config.DashDurationSeconds, config.DescriptionCharLimit, config.NameCharLimit,
	)
	c.Home()

	inputPanel.GetPrimitive().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

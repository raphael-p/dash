package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/tui/controller"
	"github.com/raphael-p/dash/internal/tui/panels"
	"github.com/raphael-p/dash/pkg/logger"
	"github.com/raphael-p/dash/pkg/tviewcomponents/commandline"
	"github.com/rivo/tview"
)

func startTUI() {
	app := tview.NewApplication()
	displayPanel := panels.NewDisplayPanel()
	infoPanel := panels.NewInfoPanel()
	navPanel := panels.NewNavigationPanel()
	commandLine := commandline.New(tcell.ColorLimeGreen)

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(navPanel.GetPrimitive(), 0, 3, false).
		AddItem(infoPanel.GetPrimitive(), 0, 1, false)
	upperFlex := tview.NewFlex().
		AddItem(displayPanel.GetPrimitive(), 0, 2, false).
		AddItem(rightFlex, 0, 1, false)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(upperFlex, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(commandLine.GetPrimitive(), 3, 1, true)

	c := controller.NewController(
		app, commandLine, navPanel, infoPanel, displayPanel,
		config.DashDurationSeconds, config.DescriptionCharLimit, config.NameCharLimit,
	)
	c.Home()

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			displayPanel.ResetPagination()
			c.Home()
			return nil
		}
		return event
	})

	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

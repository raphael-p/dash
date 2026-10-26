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
	displayPanel := panels.NewDisplayPanel(tcell.ColorGray)
	infoPanel := panels.NewInfoPanel()
	navPanel := panels.NewNavigationPanel(tcell.ColorLimeGreen)
	commandLine := commandline.New(tcell.ColorLimeGreen, tcell.ColorGray)

	minHeight, minWidth := 25, 80
	display, nav, command := displayPanel.GetPrimitive(), navPanel.GetPrimitive(), commandLine.GetPrimitive()
	spacer := tview.NewBox()

	// minimal display
	grid := tview.NewGrid().SetRows(-2, 5, 1, -1).SetColumns(-3, -1).
		AddItem(display, 0, 0, 3, 2, 0, 0, false).
		AddItem(command, 3, 0, 1, 2, 0, 0, true)

	// short display
	grid.AddItem(display, 0, 0, 3, 1, 0, minWidth, false).
		AddItem(nav, 0, 1, 3, 1, 0, minWidth, false)

	// narrow display
	grid.AddItem(display, 0, 0, 2, 2, minHeight, 0, false)

	// regular display, added last so it wins ties with the short display
	grid.AddItem(display, 0, 0, 2, 1, minHeight, minWidth, false).
		AddItem(nav, 0, 1, 1, 1, minHeight, minWidth, false).
		AddItem(infoPanel.GetPrimitive(), 1, 1, 1, 1, minHeight, minWidth, false).
		AddItem(spacer, 2, 0, 1, 2, minHeight, 0, false)

	// grid retains scroll offsets after shrinking to fit the focused command,
	// need to reset them on resize so the content returns when the terminal expands
	var lastWidth, lastHeight int
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		width, height := screen.Size()
		if width != lastWidth || height != lastHeight {
			grid.SetOffset(0, 0)
			lastWidth, lastHeight = width, height
		}
		return false
	})

	c := controller.NewController(
		app, commandLine, navPanel, infoPanel, displayPanel,
		config.DashDurationSeconds, config.DescriptionCharLimit, config.NameCharLimit,
	)
	c.Home()

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			displayPanel.ResetPagination()
			c.Home()
			return nil
		}
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

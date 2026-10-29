package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/internal/tui/components"
	"github.com/raphael-p/dash/internal/tui/panels"
	"github.com/raphael-p/dash/pkg/tviewcomponents/commandline"
	"github.com/rivo/tview"
)

type Controller struct {
	queueUpdateDraw func(func())
	stop            func()
	commandLine     *commandline.CommandLine
	navPanel        *panels.NavigationPanel
	infoPanel       *panels.InfoPanel
	displayPanel    *panels.DisplayPanel
}

func Initialise(app *tview.Application) {
	highlightColour, secondaryColour := tcell.ColorLimeGreen, tcell.ColorGray
	commandLine := commandline.New(highlightColour, secondaryColour)
	displayPanel := panels.NewDisplayPanel(secondaryColour)
	infoPanel := panels.NewInfoPanel()
	navPanel := panels.NewNavigationPanel(highlightColour, secondaryColour)

	grid := components.NewGrid(displayPanel.GetPrimitive(), infoPanel.GetPrimitive(), navPanel.GetPrimitive(), commandLine.GetPrimitive(), 25, 80)

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

	queueUpdateDraw := func(f func()) { app.QueueUpdateDraw(f) }
	c := &Controller{queueUpdateDraw, app.Stop, commandLine, navPanel, infoPanel, displayPanel}

	// global ESC shortcut that returns the user to the home screen
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			displayPanel.ResetPagination()
			c.Home()
			return nil
		}
		return event
	})

	c.Home()
	app.SetRoot(grid, true)
}

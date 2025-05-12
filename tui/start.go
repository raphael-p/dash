package tui

import (
	"github.com/raphael-p/datashard/logger"
	"github.com/rivo/tview"
)

type appController struct {
	app          *tview.Application
	inputPanel   *inputPanel
	infoPanel    *infoPanel
	displayPanel *displayPanel
}

func Start() {
	app := tview.NewApplication()

	displayPanel := makeDisplayPanel()
	infoPanel := makeInfoPanel()
	inputPanel := makeInputPanel()

	ac := &appController{app, inputPanel, infoPanel, displayPanel}

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputPanel.panel, 0, 2, true).
		AddItem(infoPanel.panel, 0, 1, false)
	flex := tview.NewFlex().
		AddItem(displayPanel.panel, 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	ac.navigateToHomeFunc(true)()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

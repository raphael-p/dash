package tui

import (
	"github.com/raphael-p/datashard/logger"
	"github.com/rivo/tview"
)

type appController struct {
	app          *tview.Application
	inputPanel   tview.Primitive
	logPanel     *tview.TextView
	displayPanel *tview.TextView
	rightFlex    *tview.Flex
}

func Start() {
	app := tview.NewApplication()

	displayPanel := tview.NewTextView().
		SetWordWrap(true)
	displayPanel.SetBorder(true).SetTitle("Tasks")

	logPanel := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	logPanel.SetBorder(true)

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	ac := &appController{app, nil, logPanel, displayPanel, rightFlex}

	flex := tview.NewFlex().
		AddItem(displayPanel, 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	ac.refreshTasks()
	ac.navigateToHomeFunc()()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

package tui

import (
	"github.com/raphael-p/datashard/logger"
	"github.com/rivo/tview"
)

type viewController struct {
	app          *tview.Application
	inputPanel   *tview.Form
	logPanel     *tview.TextView
	displayPanel *tview.TextView
}

func Home() {
	app := tview.NewApplication()

	displayPanel := tview.NewTextView().
		SetWordWrap(true)
	displayPanel.SetBorder(true).SetTitle("Tasks")

	logPanel := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	logPanel.SetBorder(true)

	inputPanel := createInputPanel(app, logPanel, displayPanel)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputPanel.inputPanel, 0, 2, true).
		AddItem(logPanel, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(displayPanel, 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	inputPanel.refreshTasks()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

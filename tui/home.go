package tui

import (
	"github.com/raphael-p/datashard/logger"
	"github.com/rivo/tview"
)

type viewController struct {
	app          *tview.Application
	inputPanel   tview.Primitive
	logPanel     *tview.TextView
	displayPanel *tview.TextView
	rightFlex    *tview.Flex
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

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	vc := &viewController{app, nil, logPanel, displayPanel, rightFlex}

	flex := tview.NewFlex().
		AddItem(displayPanel, 0, 2, false).
		AddItem(rightFlex, 0, 1, true)

	vc.refreshTasks()
	vc.navigateToHome()()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		logger.Fatal(err.Error())
	}
}

func (vc *viewController) switchInputPanel(newPrimitive tview.Primitive) {
	oldPrimitive := vc.inputPanel
	if oldPrimitive != nil {
		vc.rightFlex.RemoveItem(oldPrimitive)
	}
	vc.rightFlex.Clear().AddItem(newPrimitive, 0, 2, true)
	vc.rightFlex.AddItem(vc.logPanel, 0, 1, false)
	vc.app.SetFocus(newPrimitive)
	vc.inputPanel = newPrimitive
}

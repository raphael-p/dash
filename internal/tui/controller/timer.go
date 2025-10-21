package controller

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const refreshInterval = 2 * time.Second

func refresh(app *tview.Application, setText func()) {
	for {
		app.QueueUpdateDraw(func() {
			setText()
		})
		time.Sleep(refreshInterval)
	}
}

func (c *Controller) startTimer() {
	timerBox := tview.NewTextView()
	timerBox.SetTextColor(tcell.ColorLimeGreen)
	timerBox.SetBorderPadding(1, 0, 2, 2)

	setText := func() {
		timerBox.SetText(fmt.Sprintf("%s minutes remaining", c.inputPanel.GetMinutesRemaining()))
	}

	c.inputPanel.RestartTimer()
	go refresh(c.app, setText)

	hotkeys := tview.NewTextView().SetDynamicColors(true)
	hotkeys.SetBorderPadding(0, 1, 2, 2)
	hotkeys.SetText(fmt.Sprintf(
		`([%[1]s::b]a[-:-:-]) add a new task
([%[1]s::b]r[-:-:-]) restart timer
([%[1]s::b]c[-:-:-]) mark task as completed
([%[1]s::b]q[-:-:-]) quit`,
		tcell.ColorLimeGreen))

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(timerBox, 3, 0, false).
		AddItem(hotkeys, 0, 1, false)

	c.inputPanel.Set("Dash!", layout)
}

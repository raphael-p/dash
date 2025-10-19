package controller

import (
	"github.com/raphael-p/datashard/internal/tui/panels"
	"github.com/rivo/tview"
)

type Controller struct {
	app          *tview.Application
	inputPanel   *panels.InputPanel
	infoPanel    *panels.InfoPanel
	displayPanel *panels.DisplayPanel
}

func NewController(
	app *tview.Application,
	inputPanel *panels.InputPanel,
	infoPanel *panels.InfoPanel,
	displayPanel *panels.DisplayPanel,
) *Controller {
	return &Controller{app, inputPanel, infoPanel, displayPanel}
}
